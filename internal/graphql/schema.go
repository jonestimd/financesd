package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/jonestimd/financesd/internal/model"
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
	detailSchema.AddFieldConfig("relatedDetail", &graphql.Field{
		Type: detailSchema,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			detailsMap := DetailsCacheKey.getCache(p.Context)
			if parent, ok := p.Source.(model.TransactionDetail); ok {
				return detailsMap[*parent.RelatedDetailID], nil
			}
			return nil, nil
		},
	})
	detailSchema.AddFieldConfig("transaction", &graphql.Field{
		Type: transactionSchema,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			transactionsMap := TransactionsCacheKey.getCache(p.Context)
			if parent, ok := p.Source.(model.TransactionDetail); ok {
				return transactionsMap[parent.TransactionID], nil
			}
			return nil, nil
		},
	})

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: queries}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	return graphql.NewSchema(schemaConfig)
}

// type myExtension struct{}

// func (e *myExtension) Init(ctx context.Context, p *graphql.Params) context.Context {
// 	return ctx
// }

// func (e *myExtension) Name() string {
// 	return "My experimental extension"
// }
