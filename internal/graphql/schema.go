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
	detailSchema.AddFieldConfig("relatedDetail", &graphql.Field{Type: detailSchema})
	detailSchema.AddFieldConfig("transaction", &graphql.Field{Type: transactionSchema})

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: queries}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	return graphql.NewSchema(schemaConfig)
}
