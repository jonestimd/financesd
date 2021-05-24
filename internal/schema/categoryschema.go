package schema

import (
	"database/sql"

	"github.com/graphql-go/graphql"
)

var categorySchema = graphql.NewObject(graphql.ObjectConfig{
	Name:        "category",
	Description: "the type of a transaction",
	Fields: addAudit(graphql.Fields{
		"id":               &graphql.Field{Type: graphql.Int},
		"code":             &graphql.Field{Type: graphql.String},
		"description":      &graphql.Field{Type: graphql.String},
		"amountType":       &graphql.Field{Type: graphql.String},
		"parentId":         &graphql.Field{Type: graphql.Int},
		"security":         &graphql.Field{Type: yesNoType},
		"income":           &graphql.Field{Type: yesNoType},
		"assetExchange":    &graphql.Field{Type: yesNoType},
		"transactionCount": &graphql.Field{Type: graphql.Int},
	}),
})

var categoryQueryFields = &graphql.Field{
	Type: graphql.NewList(categorySchema),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		tx := p.Context.Value(DbContextKey).(*sql.Tx)
		return getAllCategories(tx), nil
	},
}
