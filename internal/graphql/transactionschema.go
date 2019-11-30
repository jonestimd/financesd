package graphql

import (
	"errors"
	"strconv"

	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	"github.com/jonestimd/financesd/internal/model"
)

var detailSchema = graphql.NewObject(graphql.ObjectConfig{
	Name:        "transactionDetail",
	Description: "a detail of a financial transaction",
	Fields: addAudit(graphql.Fields{
		"id":              &graphql.Field{Type: graphql.ID},
		"categoryId":      &graphql.Field{Type: graphql.Int},
		"groupId":         &graphql.Field{Type: graphql.Int},
		"memo":            &graphql.Field{Type: graphql.String},
		"amount":          &graphql.Field{Type: graphql.Float},
		"assetQuantity":   &graphql.Field{Type: graphql.Float},
		"exchangeAssetId": &graphql.Field{Type: graphql.Int},
		// add related detail in Schema()
		// add transaction in Schema()
	}),
})

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
		db := p.Context.Value(DbContextKey).(*gorm.DB)
		prepare := emptyPrepare
		if isPathSelected(p.Info, "details") {
			prepare = preload("Details")
		}
		if id, ok := p.Args["accountId"]; ok {
			transactions := make([]*model.Transaction, 0)
			if intId, err := strconv.ParseInt(id.(string), 10, 64); err == nil {
				if err := prepare(db).Order("Date").Find(&transactions, "account_id = ?", intId).Error; err != nil {
					return nil, err
				}
				if isPathSelected(p.Info, "details", "relatedDetail") {
					relatedDetailIDs := getRelatedDetailIDs(transactions)
					relatedDetails := make([]*model.TransactionDetail, 0)
					if err := db.Find(&relatedDetails, "id in (?)", relatedDetailIDs).Error; err != nil {
						return nil, err
					}
					detailsMap := DetailsCacheKey.getCache(p.Context)
					for _, detail := range relatedDetails {
						detailsMap[detail.ID] = detail
					}
					if isPathSelected(p.Info, "details", "relatedDetails", "transaction") {
						relatedTransactionIDs := getRelatedTransactionIDs(relatedDetails)
						relatedTransactions := make([]*model.Transaction, 0)
						if err := db.Find(&relatedTransactions, "id in (?)", relatedTransactionIDs).Error; err != nil {
							return nil, err
						}
						transactionsMap := TransactionsCacheKey.getCache(p.Context)
						for _, transaction := range relatedTransactions {
							transactionsMap[transaction.ID] = transaction
						}
					}
				}
				return transactions, nil
			} else {
				return nil, err
			}
		}
		return nil, errors.New("accountId is required")
	},
}

type TransactionList []*model.Transaction

func getRelatedDetailIDs(transactions []*model.Transaction) []int {
	relatedIDs := make([]int, 0)
	for _, transaction := range transactions {
		relatedIDs = transaction.GetRelatedDetailIDs(relatedIDs)
	}
	return relatedIDs
}

func getRelatedTransactionIDs(details []*model.TransactionDetail) []int {
	relatedIDs := make([]int, 0)
	for _, detail := range details {
		relatedIDs = append(relatedIDs, detail.TransactionID)
	}
	return relatedIDs
}
