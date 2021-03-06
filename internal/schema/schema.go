package schema

import (
	"github.com/graphql-go/graphql"
)

const accountQuery = "accounts"
const companyQuery = "companies"
const updateCompaniesMutation = "updateCompanies"
const payeeQuery = "payees"

// const assetsQuery = "assets"
const securityQuery = "securities"
const categoryQuery = "categories"
const groupQuery = "groups"
const transactionQuery = "transactions"
const updateTxMutation = "updateTransactions"

var queries = graphql.Fields{
	accountQuery:     accountQueryFields,
	companyQuery:     companyQueryFields(),
	payeeQuery:       payeeQueryFields,
	securityQuery:    securityQueryFields,
	categoryQuery:    categoryQueryFields,
	groupQuery:       groupQueryFields,
	transactionQuery: transactionQueryFields,
}

var mutations = graphql.Fields{
	updateCompaniesMutation: updateCompaniesFields,
	updateTxMutation:        updateTxFields,
}

// New creates the GraphQL schema.
func New() (graphql.Schema, error) {
	schemaConfig := graphql.SchemaConfig{
		Query:    graphql.NewObject(graphql.ObjectConfig{Name: "Queries", Fields: queries}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{Name: "Mutations", Fields: mutations}),
	}
	return graphql.NewSchema(schemaConfig)
}
