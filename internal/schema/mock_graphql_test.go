package schema

import (
	"context"
	"database/sql"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
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

type resolveParamsBuilder struct {
	graphql.ResolveParams
}

func newResolveInfo(queryName string, querySelection ...ast.Selection) graphql.ResolveInfo {
	schema, _ := New()
	queryKey := "the query"
	selections := []ast.Selection{
		newField(queryKey, queryName, querySelection...),
	}
	operation := mockDefinition{selections: ast.NewSelectionSet(&ast.SelectionSet{Selections: selections})}
	return graphql.ResolveInfo{
		Schema:    schema,
		Operation: &operation,
		Path:      &graphql.ResponsePath{Key: queryKey},
		RootValue: make(map[string]interface{}),
	}
}

func newResolveParams(tx *sql.Tx, queryName string, querySelection ...ast.Selection) *resolveParamsBuilder {
	info := newResolveInfo(queryName, querySelection...)
	ctx := context.WithValue(context.Background(), DbContextKey, tx)
	ctx = context.WithValue(ctx, UserKey, "somebody")
	return &resolveParamsBuilder{graphql.ResolveParams{Info: info, Context: ctx}}
}

func (b *resolveParamsBuilder) addArg(name string, value interface{}) *resolveParamsBuilder {
	if name != "" {
		if b.Args == nil {
			b.Args = make(map[string]interface{})
		}
		b.Args[name] = value
	}
	return b
}

func (b *resolveParamsBuilder) addArrayArg(name string, maps []map[string]interface{}, nested ...string) *resolveParamsBuilder {
	if name != "" {
		if b.Args == nil {
			b.Args = make(map[string]interface{})
		}
		b.Args[name] = asInterfaces(maps, nested...)
	}
	return b
}

// AsInterfaces converts a []map[string]interface{} to a []interface{}.
func asInterfaces(maps []map[string]interface{}, nested ...string) []interface{} {
	values := make([]interface{}, len(maps))
	for i, item := range maps {
		m := make(map[string]interface{})
		for key, value := range item {
			m[key] = value
		}
		for _, name := range nested {
			if value, ok := m[name]; ok {
				m[name] = asInterfaces(value.([]map[string]interface{}))
			}
		}
		values[i] = m
	}
	return values
}

func (b *resolveParamsBuilder) setSource(source interface{}) *resolveParamsBuilder {
	b.Source = source
	return b
}

func findSchemaField(root *graphql.Object, names ...string) *graphql.FieldDefinition {
	field := root.Fields()[names[0]]
	for _, name := range names[1:] {
		switch field.Type.(type) {
		case *graphql.List:
			field = field.Type.(*graphql.List).OfType.(*graphql.Object).Fields()[name]
		case *graphql.Object:
			field = field.Type.(*graphql.Object).Fields()[name]
		default:
			return nil
		}
	}
	return field
}
