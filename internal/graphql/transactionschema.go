package graphql

import (
	"errors"
	"strconv"

	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
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

type prepareDb func(db *gorm.DB) *gorm.DB

func emptyPrepare(db *gorm.DB) *gorm.DB {
	return db
}

func preload(association string) prepareDb {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload(association)
	}
}

var transactionQueryFields = &graphql.Field{
	Type: graphql.NewList(transactionSchema),
	Args: map[string]*graphql.ArgumentConfig{
		"accountId": &graphql.ArgumentConfig{Type: graphql.ID, Description: "account ID"},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		if id, ok := p.Args["accountId"]; ok {
			if intId, err := strconv.ParseInt(id.(string), 10, 64); err == nil {
				db := p.Context.Value(DbContextKey).(*gorm.DB)
				query := NewQuery("transaction", "t").Convert(p.Info).
					Where("%s.account_id = ?", intId).
					OrderBy("%[1]s.date, %[1]s.id")
				return query.Execute(db)
			} else {
				return nil, err
			}
		}
		return nil, errors.New("accountId is required")
	},
}
