package graphql

import (
	"database/sql"
	"errors"

	"github.com/graphql-go/graphql"
)

var detailSchema = graphql.NewObject(graphql.ObjectConfig{
	Name:        "transactionDetail",
	Description: "a detail of a financial transaction",
	Fields: addAudit(graphql.Fields{
		"id":                    &graphql.Field{Type: graphql.ID},
		"transactionCategoryId": &graphql.Field{Type: graphql.Int},
		"transactionGroupId":    &graphql.Field{Type: graphql.Int},
		"memo":                  &graphql.Field{Type: graphql.String},
		"amount":                &graphql.Field{Type: graphql.Float},
		"assetQuantity":         &graphql.Field{Type: graphql.Float},
		"exchangeAssetId":       &graphql.Field{Type: graphql.Int},
		// add related detail in Schema()
		// add transaction in Schema()
	}),
})

var transactionSchema = graphql.NewObject(graphql.ObjectConfig{
	Name:        "transaction",
	Description: "a financial transaction",
	Fields: addAudit(graphql.Fields{
		"id":              &graphql.Field{Type: graphql.ID},
		"date":            &graphql.Field{Type: graphql.String},
		"memo":            &graphql.Field{Type: graphql.String},
		"referenceNumber": &graphql.Field{Type: graphql.String},
		"cleared":         &graphql.Field{Type: yesNoType},
		"accountId":       &graphql.Field{Type: graphql.Int},
		"payeeId":         &graphql.Field{Type: graphql.Int},
		"securityId":      &graphql.Field{Type: graphql.Int},
		"details":         &graphql.Field{Type: graphql.NewList(detailSchema)},
	}),
})

var transactionQueryFields = &graphql.Field{
	Type: graphql.NewList(transactionSchema),
	Args: map[string]*graphql.ArgumentConfig{
		"accountId": {Type: graphql.ID, Description: "account ID"},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		if _, ok := p.Args["accountId"]; ok {
			db := p.Context.Value(DbContextKey).(*sql.Tx)
			return newQuery("transaction", "t").SelectFields(p.Info).Filter(p.Args).OrderBy("%[1]s.date, %[1]s.id").Execute(db)
		}
		return nil, errors.New("accountId is required")
	},
}
