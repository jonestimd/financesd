package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
)

var categorySchema = graphql.NewObject(graphql.ObjectConfig{
	Name:        "category",
	Description: "the type of a transaction",
	Fields: addAudit(graphql.Fields{
		"id":          &graphql.Field{Type: graphql.ID},
		"code":        &graphql.Field{Type: graphql.String},
		"description": &graphql.Field{Type: graphql.String},
		"amountType":  &graphql.Field{Type: graphql.String},
		"parentId":    &graphql.Field{Type: graphql.Int},
		"security":    &graphql.Field{Type: yesNoType},
		"income":      &graphql.Field{Type: yesNoType},
	}),
})

var categoryQueryFields = &graphql.Field{
	Type: graphql.NewList(categorySchema),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		db := p.Context.Value(DbContextKey).(*gorm.DB)
		return NewQuery("transaction_category", "c").Convert(p.Info).Execute(db)
	},
}
