package graphql

import "github.com/graphql-go/graphql/language/ast"

type mockDefinition struct {
	selections *ast.SelectionSet
}

func (o *mockDefinition) GetOperation() string {
	return "mock operation"
}

func (o *mockDefinition) GetVariableDefinitions() []*ast.VariableDefinition {
	return nil
}

func (o *mockDefinition) GetSelectionSet() *ast.SelectionSet {
	return o.selections
}

func (o *mockDefinition) GetKind() string {
	return ""
}

func (o *mockDefinition) GetLoc() *ast.Location {
	return nil
}

func newField(alias string, name string, selections ...ast.Selection) *ast.Field {
	var astAlias *ast.Name
	if alias == "" {
		astAlias = nil
	} else {
		astAlias = &ast.Name{Value: alias}
	}
	return ast.NewField(&ast.Field{
		Alias:        astAlias,
		Name:         &ast.Name{Value: name},
		SelectionSet: ast.NewSelectionSet(&ast.SelectionSet{Selections: selections}),
	})
}
