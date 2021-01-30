package schema

import (
	"database/sql"

	"github.com/graphql-go/graphql"
)

var groupSchema = graphql.NewObject(graphql.ObjectConfig{
	Name:        "group",
	Description: "alternate categorization for transaction details",
	Fields: addAudit(graphql.Fields{
		"id":               &graphql.Field{Type: graphql.ID},
		"name":             &graphql.Field{Type: graphql.String},
		"description":      &graphql.Field{Type: graphql.String},
		"transactionCount": &graphql.Field{Type: graphql.Int},
	}),
})

var groupQueryFields = &graphql.Field{
	Type: graphql.NewList(groupSchema),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		tx := p.Context.Value(DbContextKey).(*sql.Tx)
		return getAllGroups(tx)
	},
}
