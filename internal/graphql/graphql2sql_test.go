package graphql

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/graphql-go/graphql/language/ast"
)

var whitespace = regexp.MustCompile("[ \n\t]+")

func compareColumns(expected []string, actual []string, t *testing.T) {
	if len(expected) != len(actual) {
		t.Errorf("Expected %d columns, got %d", len(expected), len(actual))
	}
	normalized := make([]string, len(expected))
	for i := range expected {
		normalized[i] = whitespace.ReplaceAllLiteralString(actual[i], " ")
	}
	if !reflect.DeepEqual(expected, normalized) {
		t.Errorf("Expected %v, got %v", expected, normalized)
	}
}

func Test_graphql2sql_SelectFields_columns(t *testing.T) {
	alias := "a"
	tests := []struct {
		name           string
		fieldName      string
		expectedColumn string
	}{
		{"converts column name", "companyId", fmt.Sprintf("\"companyId\", %s.company_id", alias)},
		{"appends derived column", "transactionCount", fmt.Sprintf("\"transactionCount\", "+accountFieldSQL["transactionCount"], alias)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			q := newQuery("account", alias)
			resolveInfo := newResolveInfo(accountQuery, newField("", test.fieldName))

			q.SelectFields(resolveInfo, accountFieldSQL)

			compareColumns([]string{test.expectedColumn}, q.columns, t)
		})
	}
}

func Test_graphql2sql_SelectFields_subquery(t *testing.T) {
	tests := []struct {
		name           string
		graphqlQuery   string
		fieldName      string
		selectedField  *ast.Field
		expectedColumn string
	}{
		{"appends to-one subquery", accountQuery, "company", newField("", "name"),
			`"company", ( select json_object("name", c0.name) from company c0 where t.company_id = c0.id)`},
		{"appends list subquery", companyQuery, "accounts", newField("", "name"),
			`"accounts", ( select json_arrayagg(json_object("name", a0.name)) from account a0 where t.id = a0.table_id)`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			q := newQuery("table", "t")
			resolveInfo := newResolveInfo(test.graphqlQuery, newField("", test.fieldName, ast.Selection(test.selectedField)))

			q.SelectFields(resolveInfo, accountFieldSQL)

			compareColumns([]string{test.expectedColumn}, q.columns, t)
		})
	}
}

func Test_getAlias(t *testing.T) {
	tests := []struct {
		name  string
		table string
		alias string
	}{
		{"uses table initial", "table", "t"},
		{"uses table initials", "big_table_name", "btn"},
		{"ignores trailing _", "big_table_", "bt"},
		{"ignores multiple _s", "big___table", "bt"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := getAlias(test.table)

			if result != test.alias {
				t.Errorf("Expected %s, got %s", test.alias, result)
			}
		})
	}
}
