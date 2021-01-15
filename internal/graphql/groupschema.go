package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
)

var groupSchema = graphql.NewObject(graphql.ObjectConfig{
	Name:        "group",
	Description: "alternate categorization for transaction details",
	Fields: addAudit(graphql.Fields{
		"id":          &graphql.Field{Type: graphql.ID},
		"name":        &graphql.Field{Type: graphql.String},
		"description": &graphql.Field{Type: graphql.String},
	}),
})

var groupQueryFields = &graphql.Field{
	Type: graphql.NewList(groupSchema),
	Args: map[string]*graphql.ArgumentConfig{
		"id":   {Type: graphql.ID, Description: "group ID"},
		"name": {Type: graphql.String, Description: "unique group name"},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		db := p.Context.Value(DbContextKey).(*gorm.DB).CommonDB()
		return newQuery("tx_group", "g").SelectFields(p.Info).Filter(p.Args).Execute(db)
	},
}
