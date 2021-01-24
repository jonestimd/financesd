package graphql

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/graphql-go/graphql"
	"github.com/jonestimd/financesd/internal/model"
)

func getDetailSchema(name string, relatedField string, fieldType graphql.Output, resolve graphql.FieldResolveFn) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name:        name,
		Description: "a detail of a financial transaction",
		Fields: addAudit(graphql.Fields{
			"id":                    &graphql.Field{Type: graphql.ID},
			"transactionCategoryId": &graphql.Field{Type: graphql.Int},
			"transactionGroupId":    &graphql.Field{Type: graphql.Int},
			"memo":                  &graphql.Field{Type: graphql.String},
			"amount":                &graphql.Field{Type: graphql.Float},
			"assetQuantity":         &graphql.Field{Type: graphql.Float},
			"exchangeAssetId":       &graphql.Field{Type: graphql.Int},
			relatedField:            &graphql.Field{Type: fieldType, Resolve: resolve},
		}),
	})
}

func getTxSchemaConfig(name string) graphql.ObjectConfig {
	return graphql.ObjectConfig{
		Description: "a financial transaction",
		Name:        name,
		Fields: addAudit(graphql.Fields{
			"id":              &graphql.Field{Type: graphql.ID},
			"date":            &graphql.Field{Type: graphql.String},
			"memo":            &graphql.Field{Type: graphql.String},
			"referenceNumber": &graphql.Field{Type: graphql.String},
			"cleared":         &graphql.Field{Type: yesNoType},
			"accountId":       &graphql.Field{Type: graphql.Int},
			"payeeId":         &graphql.Field{Type: graphql.Int},
			"securityId":      &graphql.Field{Type: graphql.Int},
		}),
	}
}

func getTxSchema() *graphql.Object {
	relatedTxSchema := graphql.NewObject(getTxSchemaConfig("relatedTransaction"))
	relatedDetailSchema := getDetailSchema("relatedTransactionDetail", "transaction", relatedTxSchema, resolveRelatedTransaction)
	detailSchema := getDetailSchema("transactionDetail", "relatedDetail", relatedDetailSchema, resolveRelatedDetail)
	txSchema := graphql.NewObject(getTxSchemaConfig("transaction"))
	txSchema.AddFieldConfig("details", &graphql.Field{Type: graphql.NewList(detailSchema), Resolve: resolveDetails})
	return txSchema
}

var transactionQueryFields = &graphql.Field{
	Type: graphql.NewList(getTxSchema()),
	Args: map[string]*graphql.ArgumentConfig{
		"accountId": {Type: graphql.NewNonNull(graphql.ID), Description: "account ID"},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		tx := p.Context.Value(DbContextKey).(*sql.Tx)
		accountID, err := strconv.ParseInt(p.Args["accountId"].(string), 10, 64)
		if err != nil {
			return nil, err
		}
		return getAccountTransactions(tx, accountID)
	},
}

type txModel interface {
	GetDetails(tx *sql.Tx) ([]*model.TransactionDetail, error)
}

func resolveDetails(p graphql.ResolveParams) (interface{}, error) {
	if transaction, ok := p.Source.(txModel); ok {
		tx := p.Context.Value(DbContextKey).(*sql.Tx)
		return transaction.GetDetails(tx)
	}
	return nil, errors.New("invalid source")
}

type detailModel interface {
	GetRelatedDetail(tx *sql.Tx) (*model.TransactionDetail, error)
	GetRelatedTransaction(tx *sql.Tx) (*model.Transaction, error)
}

func resolveRelatedDetail(p graphql.ResolveParams) (interface{}, error) {
	if detail, ok := p.Source.(detailModel); ok {
		tx := p.Context.Value(DbContextKey).(*sql.Tx)
		return detail.GetRelatedDetail(tx)
	}
	return nil, errors.New("invalid source")
}

func resolveRelatedTransaction(p graphql.ResolveParams) (interface{}, error) {
	if detail, ok := p.Source.(detailModel); ok {
		tx := p.Context.Value(DbContextKey).(*sql.Tx)
		return detail.GetRelatedTransaction(tx)
	}
	return nil, errors.New("invalid source")
}
