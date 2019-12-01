package graphql

import (
	"strconv"

	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
)

// Schema
var companySchema = graphql.NewObject(graphql.ObjectConfig{
	Name:        "company",
	Description: "a financial company",
	Fields: addAudit(graphql.Fields{
		"id":   &graphql.Field{Type: graphql.ID},
		"name": &graphql.Field{Type: graphql.String},
		// add later in Schema(): "accounts": &graphql.Field{Type: graphql.NewList(accountSchema())},
	}),
})

var companyQueryFields = &graphql.Field{
	Type: graphql.NewList(companySchema),
	Args: map[string]*graphql.ArgumentConfig{
		"id":   &graphql.ArgumentConfig{Type: graphql.ID, Description: "company ID"},
		"name": &graphql.ArgumentConfig{Type: graphql.String, Description: "unique company name"},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		db := p.Context.Value(DbContextKey).(*gorm.DB)
		query := NewQuery("company", "c").SelectFields(p.Info)
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
