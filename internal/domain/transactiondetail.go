package domain

import (
	"database/sql"
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/jonestimd/financesd/internal/database"
	"github.com/jonestimd/financesd/internal/database/table"
)

// TransactionDetail represents a line item of a financial transaction.
type TransactionDetail struct {
	txSource *transactionSource
	*table.TransactionDetail
}

func NewTransactionDetail(id int64, txID int64) *TransactionDetail {
	return &TransactionDetail{TransactionDetail: &table.TransactionDetail{ID: id, TransactionID: txID}}
}

func (d *TransactionDetail) Resolve(p graphql.ResolveParams) (interface{}, error) {
	return defaultResolveFn(replaceSource(p, d.TransactionDetail))
}

// GetRelatedDetail returns the related detail.
func (d *TransactionDetail) GetRelatedDetail(tx *sql.Tx) *TransactionDetail {
	if d.RelatedDetailID == nil {
		return nil
	}
	d.txSource.loadRelatedDetails(tx)
	return d.txSource.relatedDetailsByID[*d.RelatedDetailID]
}

// GetRelatedTransaction returns the transaction for a related detail.
func (d *TransactionDetail) GetRelatedTransaction(tx *sql.Tx) *Transaction {
	d.txSource.loadRelatedTransactions(tx)
	return d.txSource.relatedTxByID[d.TransactionID]
}

var updateTxDetails = func(tx *sql.Tx, txId int64, details []map[string]interface{}, user string) {
	deleteIDs := make([]*database.VersionID, 0)
	for _, detail := range details {
		values := database.InputObject(detail)
		id := values.GetVersionID()
		if id != nil {
			if len(detail) == 2 {
				deleteIDs = append(deleteIDs, id)
			} else {
				amount, setAmount := values.GetFloat("amount")
				transferAccountId, setTransfer := values.GetInt("transferAccountId")
				categoryId, setCategory := values.GetInt("transactionCategoryId")
				if setTransfer && setCategory {
					panic(errors.New("cannot specify both transferAccountId and transactionCategoryId"))
				}
				if setTransfer {
					setCategory = true
					categoryId = nil
				}
				updateDetail(tx, id.ID, id.Version, setCategory, categoryId, values, user)
				if setCategory && transferAccountId == nil {
					deleteTransfer(tx, id.ID)
				} else if setTransfer {
					addOrUpdateTransfer(tx, id.ID, transferAccountId, user)
				}
				if setAmount && !setCategory {
					setTransferAmount(tx, id.ID, amount, user)
				}
			}
		} else {
			if amount, ok := values.GetFloat("amount"); ok && amount != nil {
				insertDetail(tx, amount, values, user)
			} else {
				panic(errors.New("amount is required to add a transaction detail"))
			}
		}
	}
	if len(deleteIDs) > 0 {
		deleteDetails(tx, deleteIDs)
	}
}
