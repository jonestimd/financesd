package graphql

import (
	"strconv"

	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
)

// Schema
var accountSchema = graphql.NewObject(graphql.ObjectConfig{
	Name:        "account",
	Description: "a financial account",
	Fields: addAudit(graphql.Fields{
		"id":          &graphql.Field{Type: graphql.ID},
		"company":     &graphql.Field{Type: companySchema},
		"name":        &graphql.Field{Type: graphql.String},
		"description": &graphql.Field{Type: graphql.String},
		"accountNo":   &graphql.Field{Type: graphql.String},
		"type":        &graphql.Field{Type: graphql.String}, // TODO enum?
		"closed":      &graphql.Field{Type: yesNoType},
		"currencyId":  &graphql.Field{Type: graphql.Int},
	}),
})

var accountQueryFields = &graphql.Field{
	Type: graphql.NewList(accountSchema),
	Args: map[string]*graphql.ArgumentConfig{
		"id":   &graphql.ArgumentConfig{Type: graphql.ID, Description: "account ID"},
		"name": &graphql.ArgumentConfig{Type: graphql.String, Description: "unique account name"},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		query := NewQuery("account", "a").SelectFields(p.Info)
		db := p.Context.Value(DbContextKey).(*gorm.DB)
		if id, ok := p.Args["id"]; ok {
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
