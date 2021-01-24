package model

import (
	"database/sql"
	"reflect"
)

// TransactionDetail represents a line item of a financial transaction.
type TransactionDetail struct {
	ID                    int64
	TransactionID         int64
	Transaction           *Transaction
	TransactionCategoryID *int64
	TransactionGroupID    *int64
	Memo                  *string
	Amount                float64
	AssetQuantity         *float64
	ExchangeAssetID       *int64
	RelatedDetailID       *int64
	RelatedDetail         *TransactionDetail
	Version               int
	accountTransactions   *AccountTransactions
	Audited
}

func (d *TransactionDetail) ptrTo(column string) interface{} {
	switch column {
	case "id":
		return &d.ID
	case "transaction_id":
		return &d.TransactionID
	case "transaction_category_id":
		return &d.TransactionCategoryID
	case "transaction_group_id":
		return &d.TransactionGroupID
	case "memo":
		return &d.Memo
	case "amount":
		return &d.Amount
	case "asset_quantity":
		return &d.AssetQuantity
	case "exchange_asset_id":
		return &d.ExchangeAssetID
	case "related_detail_id":
		return &d.RelatedDetailID
	case "version":
		return &d.Version
	}
	return d.Audited.ptrToAudit(column)
}

var detailType = reflect.TypeOf(TransactionDetail{})

// GetRelatedDetail returns the related detail.
func (d *TransactionDetail) GetRelatedDetail(tx *sql.Tx) (*TransactionDetail, error) {
	if d.RelatedDetailID == nil {
		return nil, nil
	}
	if d.accountTransactions.err != nil {
		return nil, d.accountTransactions.err
	}
	if d.accountTransactions.relatedDetailsByID == nil && d.accountTransactions.getRelatedDetails(tx) != nil {
		return nil, d.accountTransactions.err
	}
	return d.accountTransactions.relatedDetailsByID[*d.RelatedDetailID], nil
}

// GetRelatedTransaction returns the transaction for a related detail.
func (d *TransactionDetail) GetRelatedTransaction(tx *sql.Tx) (*Transaction, error) {
	if d.accountTransactions.err != nil {
		return nil, d.accountTransactions.err
	}
	if d.accountTransactions.relatedTxByID == nil && d.accountTransactions.getRelatedTransactions(tx) != nil {
		return nil, d.accountTransactions.err
	}
	return d.accountTransactions.relatedTxByID[d.TransactionID], nil
}
