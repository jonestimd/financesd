package graphql

import (
	"errors"
	"strconv"

	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	"github.com/jonestimd/financesd/internal/model"
)

var transactionSchema = graphql.NewObject(graphql.ObjectConfig{
	Name:        "transaction",
	Description: "a financial transaction",
	Fields: addAudit(graphql.Fields{
		"id":              &graphql.Field{Type: graphql.ID},
		"date":            &graphql.Field{Type: dateType},
		"memo":            &graphql.Field{Type: graphql.String},
		"referenceNumber": &graphql.Field{Type: graphql.String},
		"cleared":         &graphql.Field{Type: yesNoType},
		"accountId":       &graphql.Field{Type: graphql.Int},
		"payeeId":         &graphql.Field{Type: graphql.Int},
		"securityId":      &graphql.Field{Type: graphql.Int},
	}),
})

var transactionQueryFields = &graphql.Field{
	Type: graphql.NewList(transactionSchema),
	Args: map[string]*graphql.ArgumentConfig{
		"accountId": &graphql.ArgumentConfig{Type: graphql.ID, Description: "account ID"},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		db := p.Context.Value(DbContextKey).(*gorm.DB)
		// if isSelected(accountQuery, p.Info, "company") {
		// 	db = db.Preload("Company") // TODO make sure it only applies to account query and doesn't effect TX
		// }
		if id, ok := p.Args["accountId"]; ok {
			transactions := make([]*model.Transaction, 0)
			if intId, err := strconv.ParseInt(id.(string), 10, 64); err == nil {
				return transactions, db.Find(&transactions, "account_id = ?", intId).Error
			} else {
				return nil, err
			}
		}
		return nil, errors.New("accountId is required")
	},
}
