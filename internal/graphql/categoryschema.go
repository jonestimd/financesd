package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
)

var categorySchema = graphql.NewObject(graphql.ObjectConfig{
	Name:        "category",
	Description: "the type of a transaction",
	Fields: addAudit(graphql.Fields{
		"id":               &graphql.Field{Type: graphql.ID},
		"code":             &graphql.Field{Type: graphql.String},
		"description":      &graphql.Field{Type: graphql.String},
		"amountType":       &graphql.Field{Type: graphql.String},
		"parentId":         &graphql.Field{Type: graphql.Int},
		"security":         &graphql.Field{Type: yesNoType},
		"income":           &graphql.Field{Type: yesNoType},
		"transactionCount": &graphql.Field{Type: graphql.Int},
	}),
})

var categoryFieldSql = map[string]string{
	"transactionCount": "(select count(distinct transaction_id) from transaction_detail where transaction_category_id = %s.id)",
}

var categoryQueryFields = &graphql.Field{
	Type: graphql.NewList(categorySchema),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		db := p.Context.Value(DbContextKey).(*gorm.DB)
		return NewQuery("transaction_category", "c").SelectFields(p.Info, categoryFieldSql).Execute(db)
	},
}
