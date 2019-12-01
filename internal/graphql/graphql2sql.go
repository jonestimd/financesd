package graphql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/jinzhu/gorm"
)

type fieldConfig struct {
	table   string
	alias   string
	where   string
	aggFunc *string
}

var jsonArrayagg = "json_arrayagg"

var fieldConfigMap = map[string]fieldConfig{
	"company":       fieldConfig{table: "company", alias: "c", where: "%s.company_id = %%s.id"},
	"accounts":      fieldConfig{table: "account", alias: "a", where: "%s.id = %%s.company_id", aggFunc: &jsonArrayagg},
	"details":       fieldConfig{table: "transaction_detail", alias: "td", where: "%s.id = %%s.transaction_id", aggFunc: &jsonArrayagg},
	"relatedDetail": fieldConfig{table: "transaction_detail", alias: "rd", where: "%s.related_detail_id = %%s.id"},
	"transaction":   fieldConfig{table: "transaction", alias: "rt", where: "%s.transaction_id = %%s.id"},
}

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

// NewQuery creates a SQL builder for retrieving JSON data.
func NewQuery(table string, alias string) *sqlData {
	return &sqlData{table: table, alias: alias, columns: make([]string, 0)}
}

func subQuery(parent *sqlData, config fieldConfig) *sqlData {
	q := NewQuery(config.table, fmt.Sprintf("%s%d", config.alias, parent.depth))
	q.depth = parent.depth + 1
	q.aggFunc = config.aggFunc
	return q
}

func (q *sqlData) Convert(info graphql.ResolveInfo) *sqlData { // TODO rename to Select()?
	return q.Append(findQuery(info).GetSelectionSet().Selections)
}

func (q *sqlData) Append(selection []ast.Selection) *sqlData {
	for _, field := range selection {
		switch field := field.(type) {
		case *ast.Field:
			if field.GetSelectionSet() == nil {
				q.columns = append(q.columns, fmt.Sprintf("\"%s\", %s.%s", field.Name.Value, q.alias, toColumn(field.Name.Value)))
			} else {
				q.addSubQuery(field.Name.Value, fieldConfigMap[field.Name.Value], field.GetSelectionSet().Selections)
			}
		}
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

func (q *sqlData) addSubQuery(name string, config fieldConfig, selection []ast.Selection) {
	subQuery := subQuery(q, config)
	subQuery.Append(selection).Where(fmt.Sprintf(config.where, q.alias))
	q.columns = append(q.columns, fmt.Sprintf("\"%s\", (%s)", name, subQuery))
}

func (q *sqlData) Execute(db *gorm.DB) ([]interface{}, error) {
	if rows, err := db.CommonDB().Query(q.String(), q.Args...); err != nil {
		return nil, err
	} else {
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
}

func (q *sqlData) String() string {
	builder := &strings.Builder{}
	builder.WriteString("select ")
	if q.aggFunc != nil {
		builder.WriteString(*q.aggFunc)
		builder.WriteRune('(')
	}
	builder.WriteString(fmt.Sprintf("json_object(%s)", strings.Join(q.columns, ", ")))
	if q.aggFunc != nil {
		builder.WriteRune(')')
	}
	builder.WriteString(fmt.Sprintf("\nfrom %s %s", q.table, q.alias))
	if q.conditions != nil {
		builder.WriteString("\nwhere ")
		builder.WriteString(strings.Join(q.conditions, " and "))
	}
	if q.orderBy != nil {
		builder.WriteString("\norder by ")
		builder.WriteString(fmt.Sprintf(*q.orderBy, q.alias))
	}
	return builder.String()
}

func toColumn(fieldName string) string {
	re := regexp.MustCompile("[A-Z]")
	return re.ReplaceAllStringFunc(fieldName, func(match string) string {
		return "_" + strings.ToLower(match)
	})
}
