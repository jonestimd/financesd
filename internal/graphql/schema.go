package graphql

import (
	"github.com/graphql-go/graphql"
)

type ReqContextKey string

const DbContextKey = ReqContextKey("db")

var queries = graphql.Fields{
	"accounts":  accountQueryFields,
	"companies": companyQueryFields,
}

func Schema() (graphql.Schema, error) {
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: queries}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	return graphql.NewSchema(schemaConfig)
}
