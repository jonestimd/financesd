package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

type ReqContextKey string

const DbContextKey = ReqContextKey("db")
const accountQuery = "accounts"
const companyQuery = "companies"

var queries = graphql.Fields{
	accountQuery: accountQueryFields,
	companyQuery: companyQueryFields,
}

func Schema() (graphql.Schema, error) {
	companySchema.AddFieldConfig("accounts", &graphql.Field{Type: graphql.NewList(accountSchema)})
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: queries}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	return graphql.NewSchema(schemaConfig)
}

// isSelected checks if a property has been selected for a schema in the query.
func isSelected(schemaName string, info graphql.ResolveInfo, selectionName string) bool {
	for _, set := range info.Operation.GetSelectionSet().Selections {
		switch set.(type) {
		case *ast.Field:
			if set.(*ast.Field).Name.Value == schemaName {
				for _, selection := range set.GetSelectionSet().Selections {
					switch selection.(type) {
					case *ast.Field:
						if selection.(*ast.Field).Name.Value == selectionName {
							return true
						}
					}
				}
			}
		}
	}
	return false
}
