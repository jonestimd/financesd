package graphql

import (
	"context"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/jinzhu/gorm"
)

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
	var selectionSet *ast.SelectionSet
	if len(selections) == 0 {
		selectionSet = nil
	} else {
		selectionSet = ast.NewSelectionSet(&ast.SelectionSet{Selections: selections})
	}
	return ast.NewField(&ast.Field{
		Alias:        astAlias,
		Name:         &ast.Name{Value: name},
		SelectionSet: selectionSet,
	})
}

func newResolveInfo(queryName string, querySelection ...ast.Selection) graphql.ResolveInfo {
	schema, _ := Schema()
	queryKey := "the query"
	selections := []ast.Selection{
		newField(queryKey, queryName, querySelection...),
	}
	operation := mockDefinition{selections: ast.NewSelectionSet(&ast.SelectionSet{Selections: selections})}
	return graphql.ResolveInfo{Schema: schema, Operation: &operation, Path: &graphql.ResponsePath{Key: queryKey}}
}

func newResolveParams(orm *gorm.DB, queryName string, args map[string]interface{}, querySelection ...ast.Selection) graphql.ResolveParams {
	info := newResolveInfo(queryName, querySelection...)
	context := context.WithValue(context.TODO(), DbContextKey, orm)
	return graphql.ResolveParams{Info: info, Args: args, Context: context}
}
