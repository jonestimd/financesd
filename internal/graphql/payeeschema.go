package graphql

import (
	"strconv"

	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
)

// Schema
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
		query := NewQuery("payee", "p").Convert(p.Info)
		if id, ok := p.Args["id"]; ok { // TODO pass Args to query.Convert()
			if intId, err := strconv.ParseInt(id.(string), 10, 64); err == nil {
				query = query.Where("%s.id = ?", intId)
			}
		}
		if name, ok := p.Args["name"]; ok {
			query = query.Where("%s.name = ?", name.(string))
		}
		return query.Execute(db)
	},
}
