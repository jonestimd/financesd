package graphql

import (
	"reflect"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

type ReqContextKey string

const DbContextKey = ReqContextKey("db")
const accountQuery = "accounts"
const companyQuery = "companies"
const payeeQuery = "payees"
const assetsQuery = "assets"
const securityQuery = "securities"
const categoryQuery = "categories"
const transactionQuery = "transactions"

var queries = graphql.Fields{
	accountQuery:     accountQueryFields,
	companyQuery:     companyQueryFields,
	payeeQuery:       payeeQueryFields,
	securityQuery:    securityQueryFields,
	categoryQuery:    categoryQueryFields,
	transactionQuery: transactionQueryFields,
}

func Schema() (graphql.Schema, error) {
	companySchema.AddFieldConfig("accounts", &graphql.Field{Type: graphql.NewList(accountSchema)})
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: queries}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	return graphql.NewSchema(schemaConfig)
}

// isSelected checks if a property has been selected for a schema in the query.
func isSelected(schemaName string, info graphql.ResolveInfo, selectionName string) bool {
	for _, query := range info.Operation.GetSelectionSet().Selections {
		switch queryField := query.(type) {
		case *ast.Field:
			if queryField.Name.Value == schemaName {
				for _, selection := range query.GetSelectionSet().Selections {
					switch selectionField := selection.(type) {
					case *ast.Field:
						if selectionField.Name.Value == selectionName {
							return true
						}
					}
				}
			}
		}
	}
	return false
}

func nestedResolver(fieldNames ...string) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		value := reflect.ValueOf(p.Source)
		for _, fieldName := range fieldNames {
			if value.Kind() == reflect.Interface || value.Kind() == reflect.Ptr {
				value = value.Elem()
			}
			value = value.FieldByName(fieldName)
		}
		return value.Interface(), nil
	}
}
