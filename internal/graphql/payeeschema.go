package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
)

var payeeSchema = graphql.NewObject(graphql.ObjectConfig{
	Name:        "payee",
	Description: "the other party in a transaction",
	Fields: addAudit(graphql.Fields{
		"id":   &graphql.Field{Type: graphql.ID},
		"name": &graphql.Field{Type: graphql.String},
	}),
})

var payeeQueryFields = &graphql.Field{
	Type: graphql.NewList(payeeSchema),
	Args: map[string]*graphql.ArgumentConfig{
		"id":   &graphql.ArgumentConfig{Type: graphql.ID, Description: "payee ID"},
		"name": &graphql.ArgumentConfig{Type: graphql.String, Description: "unique payee name"},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		db := p.Context.Value(DbContextKey).(*gorm.DB)
		return NewQuery("payee", "p").SelectFields(p.Info).Filter(p.Args).Execute(db)
	},
}
