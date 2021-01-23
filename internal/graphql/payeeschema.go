package graphql

import (
	"database/sql"

	"github.com/graphql-go/graphql"
)

var payeeSchema = graphql.NewObject(graphql.ObjectConfig{
	Name:        "payee",
	Description: "the other party in a transaction",
	Fields: addAudit(graphql.Fields{
		"id":               &graphql.Field{Type: graphql.ID},
		"name":             &graphql.Field{Type: graphql.String},
		"transactionCount": &graphql.Field{Type: graphql.Int},
	}),
})

var payeeQueryFields = &graphql.Field{
	Type: graphql.NewList(payeeSchema),
	Args: map[string]*graphql.ArgumentConfig{
		"id":   {Type: graphql.ID, Description: "payee ID"},
		"name": {Type: graphql.String, Description: "unique payee name"},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		tx := p.Context.Value(DbContextKey).(*sql.Tx)
		return getAllPayees(tx)
	},
}
