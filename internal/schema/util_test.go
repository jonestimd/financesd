package schema

import (
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

func Test_findField(t *testing.T) {
	field := newField("", "the field")
	tests := []struct {
		description   string
		matches       bool
		expectedField *ast.Field
	}{
		{"returns nil for no match", false, nil},
		{"returns field where predicate returns true", true, field},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			predicate := func(field *ast.Field) bool {
				return test.matches
			}
			selection := make([]ast.Selection, 0)
			selection = append(selection, field)

			result := findField(selection, predicate)

			if result != test.expectedField {
				t.Errorf("Expected: %v, got: %v", test.expectedField, result)
			}
		})
	}
}

func Test_byName(t *testing.T) {
	predicate := byName("the name")

	if predicate(newField("", "not the name")) {
		t.Error("Expected predicate not to match")
	}
	if !predicate(newField("", "the name")) {
		t.Error("Expected predicate to match")
	}
}

func Test_byAlias(t *testing.T) {
	predicate := byAlias("the alias")
	tests := []struct {
		description string
		alias       string
		name        string
		match       bool
	}{
		{"returns true if alias matches", "the alias", "not the alias", true},
		{"returns true if no alias and name matches", "", "the alias", true},
		{"returns false if alias doesn't match", "not the alias", "the alias", false},
		{"returns false if no alias and name doesn't match", "", "not the alias", false},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			if predicate(newField(test.alias, test.name)) != test.match {
				t.Errorf("Expected predicte to return %v", test.match)
			}
		})
	}
}

func Test_findQuery(t *testing.T) {
	queryKey := "the query"
	selections := []ast.Selection{
		newField("not the query", ""),
		newField(queryKey, ""),
	}
	operation := mockDefinition{selections: ast.NewSelectionSet(&ast.SelectionSet{Selections: selections})}
	info := graphql.ResolveInfo{Operation: &operation, Path: &graphql.ResponsePath{Key: queryKey}}

	result := findQuery(info)

	if result != selections[1] {
		t.Errorf("Expected the graphql query field")
	}
}

type child struct {
	ID string
}
type testParam struct {
	Child child
}

func Test_nestedResolver(t *testing.T) {
	value := "xyz"
	source := testParam{
		Child: child{ID: value},
	}
	params := graphql.ResolveParams{Source: &source}
	resolver := nestedResolver("Child", "ID")

	result, err := resolver(params)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != "xyz" {
		t.Errorf("Expected: %v, got: %v", value, result)
	}
}
