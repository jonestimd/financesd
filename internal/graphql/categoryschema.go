package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	"github.com/jonestimd/financesd/internal/model"
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
		categories := make([]*model.TransactionCategory, 0)
		return categories, db.Find(&categories).Error
	},
}
