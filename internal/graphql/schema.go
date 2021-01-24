package graphql

import (
	"github.com/graphql-go/graphql"
)

const accountQuery = "accounts"
const companyQuery = "companies"
const payeeQuery = "payees"
const assetsQuery = "assets"
const securityQuery = "securities"
const categoryQuery = "categories"
const groupQuery = "groups"
const transactionQuery = "transactions"

var queries = graphql.Fields{
	accountQuery:     accountQueryFields,
	companyQuery:     companyQueryFields(),
	payeeQuery:       payeeQueryFields,
	securityQuery:    securityQueryFields,
	categoryQuery:    categoryQueryFields,
	groupQuery:       groupQueryFields,
	transactionQuery: transactionQueryFields,
}

// Schema creates the root GraphQL schema.
func Schema() (graphql.Schema, error) {
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: queries}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	return graphql.NewSchema(schemaConfig)
}
