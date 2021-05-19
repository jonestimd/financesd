package model

import (
	"database/sql"
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/jonestimd/financesd/internal/database"
)

// TransactionDetail represents a line item of a financial transaction.
type TransactionDetail struct {
	txSource *transactionSource
	*database.TransactionDetail
}

func NewTransactionDetail(id int64, txID int64) *TransactionDetail {
	return &TransactionDetail{TransactionDetail: &database.TransactionDetail{ID: id, TransactionID: txID}}
}

func (d *TransactionDetail) Resolve(p graphql.ResolveParams) (interface{}, error) {
	return defaultResolveFn(replaceSource(p, d.TransactionDetail))
}

// GetRelatedDetail returns the related detail.
func (d *TransactionDetail) GetRelatedDetail(tx *sql.Tx) (*TransactionDetail, error) {
	if d.RelatedDetailID == nil {
		return nil, nil
	}
	if d.txSource.loadRelatedDetails(tx) != nil {
		return nil, d.txSource.err
	}
	return d.txSource.relatedDetailsByID[*d.RelatedDetailID], nil
}

// GetRelatedTransaction returns the transaction for a related detail.
func (d *TransactionDetail) GetRelatedTransaction(tx *sql.Tx) (*Transaction, error) {
	if d.txSource.loadRelatedTransactions(tx) != nil {
		return nil, d.txSource.err
	}
	return d.txSource.relatedTxByID[d.TransactionID], nil
}

var updateTxDetails = func(tx *sql.Tx, txId int64, details []map[string]interface{}, user string) error {
	deleteIDs := make([]*database.VersionID, 0)
	for _, detail := range details {
		values := database.InputObject(detail)
		if id, err := values.GetVersionID(); err != nil {
			return err
		} else if id != nil {
			if len(detail) == 2 {
				deleteIDs = append(deleteIDs, id)
			} else {
				amount, setAmount := values.GetFloat("amount")
				transferAccountId, setTransfer := values.GetInt("transferAccountId")
				categoryId, setCategory := values.GetInt("transactionCategoryId")
				if setTransfer && setCategory {
					return errors.New("cannot specify both transferAccountId and transactionCategoryId")
				}
				if setTransfer {
					setCategory = true
					categoryId = nil
				}
				if err := updateDetail(tx, id.ID, id.Version, setCategory, categoryId, values, user); err != nil {
					return err
				}
				if setCategory && transferAccountId == nil {
					if err := deleteTransfer(tx, id.ID); err != nil {
						return err
					}
				} else if setTransfer {
					if err := addOrUpdateTransfer(tx, id.ID, transferAccountId, user); err != nil {
						return err
					}
				}
				if setAmount && !setCategory {
					return setTransferAmount(tx, id.ID, amount, user)
				}
			}
		} else {
			if amount, ok := values.GetFloat("amount"); ok && amount != nil {
				if err := insertDetail(tx, amount, values, user); err != nil {
					return err
				}
			} else {
				return errors.New("amount is required to add a transaction detail")
			}
		}
	}
	if len(deleteIDs) > 0 {
		return deleteDetails(tx, deleteIDs)
	}
	return nil
}
