package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
)

// Schema
var accountSchema = graphql.NewObject(graphql.ObjectConfig{
	Name:        "account",
	Description: "a financial account",
	Fields: addAudit(graphql.Fields{
		"id":               &graphql.Field{Type: graphql.ID},
		"companyId":        &graphql.Field{Type: graphql.Int},
		"company":          &graphql.Field{Type: companySchema},
		"name":             &graphql.Field{Type: graphql.String},
		"description":      &graphql.Field{Type: graphql.String},
		"accountNo":        &graphql.Field{Type: graphql.String},
		"type":             &graphql.Field{Type: graphql.String}, // TODO enum?
		"closed":           &graphql.Field{Type: yesNoType},
		"currencyId":       &graphql.Field{Type: graphql.Int},
		"transactionCount": &graphql.Field{Type: graphql.Int},
		"balance":          &graphql.Field{Type: graphql.Float},
	}),
})

var accountFieldSQL = map[string]string{
	"transactionCount": "(select count(*) from transaction where account_id = %s.id)",
	"balance": "(select sum(td.amount) " +
		"from transaction tx " +
		"join transaction_detail td on tx.id = td.transaction_id " +
		"left join transaction_category tc on td.transaction_category_id = tc.id " +
		"where tx.account_id = %s.id and coalesce(tc.amount_type, '') != 'ASSET_VALUE')",
}

var accountQueryFields = &graphql.Field{
	Type: graphql.NewList(accountSchema),
	Args: map[string]*graphql.ArgumentConfig{
		"id":   {Type: graphql.ID, Description: "account ID"},
		"name": {Type: graphql.String, Description: "unique account name"},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		db := p.Context.Value(DbContextKey).(*gorm.DB).CommonDB()
		return newQuery("account", "a").SelectFields(p.Info, accountFieldSQL).Filter(p.Args).Execute(db)
	},
}
