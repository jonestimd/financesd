package graphql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/jinzhu/gorm"
)

var jsonArrayagg = "json_arrayagg"

type sqlData struct { // TODO unique alias per subquery
	depth      int
	table      string
	alias      string
	columns    []string
	aggFunc    *string
	conditions []string
	orderBy    *string
	Args       []interface{}
}

// newQuery creates a SQL builder for retrieving JSON data.
func newQuery(table string, alias string) *sqlData {
	return &sqlData{table: table, alias: alias, columns: make([]string, 0)}
}

func (q *sqlData) SelectFields(info graphql.ResolveInfo, derivedFieldSqls ...map[string]string) *sqlData {
	gq := findQuery(info)
	fieldType := info.Schema.QueryType().Fields()[gq.Name.Value].Type
	return q.selectFields(fieldType, gq.GetSelectionSet().Selections, derivedFieldSqls...)
}

func (q *sqlData) selectFields(outType graphql.Type, selection []ast.Selection, derivedFieldSqls ...map[string]string) *sqlData {
	derivedfieldSQL := make(map[string]string)
	for _, fieldSQL := range derivedFieldSqls {
		for field, sql := range fieldSQL {
			derivedfieldSQL[field] = sql
		}
	}
	for _, field := range selection {
		switch field := field.(type) {
		case *ast.Field:
			if field.GetSelectionSet() == nil {
				if sql, ok := derivedfieldSQL[field.Name.Value]; ok {
					q.columns = append(q.columns, fmt.Sprintf("\"%s\", %s", field.Name.Value, fmt.Sprintf(sql, q.alias)))
				} else {
					q.columns = append(q.columns, fmt.Sprintf("\"%s\", %s.%s", field.Name.Value, q.alias, toSnakeCase(field.Name.Value)))
				}
			} else {
				fieldType := getFieldType(outType, field.Name.Value)
				q.addSubQuery(field.Name.Value, fieldType, field.GetSelectionSet().Selections)
			}
		}
	}
	return q
}

func getFieldType(parentType graphql.Type, name string) graphql.Type {
	switch gt := unwrapList(parentType).(type) {
	case *graphql.Object:
		return gt.Fields()[name].Type
	}
	return nil
}

func (q *sqlData) Filter(args map[string]interface{}) *sqlData {
	for name, value := range args {
		q.Where(fmt.Sprintf("%%s.%s = ?", toSnakeCase(name)), value)
	}
	return q
}

func (q *sqlData) Where(condition string, args ...interface{}) *sqlData {
	if q.conditions == nil {
		q.conditions = make([]string, 0)
	}
	q.conditions = append(q.conditions, fmt.Sprintf(condition, q.alias))
	if len(args) > 0 {
		if q.Args == nil {
			q.Args = make([]interface{}, 0)
		}
		q.Args = append(q.Args, args...)
	}
	return q
}

func (q *sqlData) OrderBy(orderBy string) *sqlData {
	q.orderBy = &orderBy
	return q
}

func (q *sqlData) addSubQuery(fieldName string, outType graphql.Type, selection []ast.Selection) {
	table := getTableName(outType)
	subQuery := newQuery(table, fmt.Sprintf("%s%d", getAlias(table), q.depth))
	subQuery.depth = q.depth + 1
	subQuery.selectFields(outType, selection)
	if isList(outType) {
		subQuery.aggFunc = &jsonArrayagg
		subQuery.Where(fmt.Sprintf("%s.id = %%s.%s_id", q.alias, q.table))
	} else {
		subQuery.Where(fmt.Sprintf("%s.%s_id = %%s.id", q.alias, toSnakeCase(fieldName)))
	}
	q.columns = append(q.columns, fmt.Sprintf("\"%s\", (%s)", fieldName, subQuery))
}

func isList(gt graphql.Type) bool {
	_, ok := gt.(*graphql.List)
	return ok
}

func unwrapList(gt graphql.Type) graphql.Type {
	if list, ok := gt.(*graphql.List); ok {
		return list.OfType
	}
	return gt
}

func getTableName(gt graphql.Type) string {
	return toSnakeCase(unwrapList(gt).Name())
}

func (q *sqlData) Execute(db gorm.SQLCommon) ([]interface{}, error) {
	query := q.String()
	if os.Getenv("SHOW_SQL") != "" {
		log.Println("SQL:", query)
	}
	rows, err := db.Query(query, q.Args...)
	if err != nil {
		return nil, err
	}
	results := make([]interface{}, 0)
	for rows.Next() {
		var jsonString sql.RawBytes
		var result interface{}
		if err := rows.Scan(&jsonString); err != nil {
			return nil, err
		}
		json.Unmarshal(jsonString, &result)
		results = append(results, result)
	}
	return results, nil
}

func (q *sqlData) String() string {
	indent := strings.Repeat("  ", q.depth)
	builder := &strings.Builder{}
	if q.depth > 0 {
		builder.WriteRune('\n')
		builder.WriteString(indent)
	}
	builder.WriteString("select ")
	if q.aggFunc != nil {
		builder.WriteString(*q.aggFunc)
		builder.WriteRune('(')
	}
	builder.WriteString(fmt.Sprintf("json_object(%s)", strings.Join(q.columns, ", ")))
	if q.aggFunc != nil {
		builder.WriteRune(')')
	}
	builder.WriteString(fmt.Sprintf("\n%sfrom %s %s", indent, q.table, q.alias))
	if q.conditions != nil {
		builder.WriteRune('\n')
		builder.WriteString(indent)
		builder.WriteString("where ")
		builder.WriteString(strings.Join(q.conditions, " and "))
	}
	if q.orderBy != nil {
		builder.WriteRune('\n')
		builder.WriteString(indent)
		builder.WriteString("order by ")
		builder.WriteString(fmt.Sprintf(*q.orderBy, q.alias))
	}
	return builder.String()
}

func getAlias(table string) string {
	alias := table[0:1]
	append := false
	for _, rune := range table[1:] {
		if rune == '_' {
			append = true
		} else if append {
			alias += string(rune)
			append = false
		}
	}
	return alias
}

func toSnakeCase(name string) string {
	re := regexp.MustCompile("[A-Z]")
	return re.ReplaceAllStringFunc(name, func(match string) string {
		return "_" + strings.ToLower(match)
	})
}
