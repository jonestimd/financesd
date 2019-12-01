package graphql

import (
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
		return NewQuery("company", "c").SelectFields(p.Info).Filter(p.Args).Execute(db)
	},
}
