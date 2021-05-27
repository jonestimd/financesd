package domain

import (
	"database/sql"
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/jonestimd/financesd/internal/database"
	"github.com/jonestimd/financesd/internal/database/table"
)

// Transaction represents a financial transaction.
type Transaction struct {
	source *transactionSource
	*table.Transaction
}

func NewTransaction(id int64) *Transaction {
	return &Transaction{Transaction: &table.Transaction{ID: id}}
}

func (t *Transaction) Resolve(p graphql.ResolveParams) (interface{}, error) {
	return defaultResolveFn(replaceSource(p, t.Transaction))
}

// GetDetails returns details for a transaction in the account.
func (t *Transaction) GetDetails(tx *sql.Tx) []*TransactionDetail {
	t.source.loadTransactionDetails(tx)
	return t.source.detailsByTxID[t.ID]
}

// GetTransactions returns all transactions for the account.
func GetTransactions(tx *sql.Tx, accountID int64) []*Transaction {
	at := &transactionSource{accountID: accountID}
	return at.setSource(getTransactions(tx, accountID))
}

// GetTransactionsByIDs returns transactions for the specified IDs.
func GetTransactionsByIDs(tx *sql.Tx, ids []int64) []*Transaction {
	source := &transactionSource{txIDs: ids}
	return source.setSource(getTransactionsByIDs(tx, ids))
}

// InsertTransactions inserts transactions.
func InsertTransactions(tx *sql.Tx, accountID int64, inserts []map[string]interface{}, user string) []int64 {
	ids := make([]int64, len(inserts))
	for i, transaction := range inserts {
		if detailsArg, ok := transaction["details"]; !ok {
			panic(errors.New("new transaction requires at least 1 detail"))
		} else {
			details := detailsArg.([]map[string]interface{})
			if len(details) < 1 {
				panic(errors.New("new transaction requires at least 1 detail"))
			}
			ids[i] = insertTransaction(tx, accountID, transaction, user)
			for _, detail := range details {
				if amount, ok := database.InputObject(detail).GetFloat("amount"); ok && amount != nil {
					insertDetail(tx, ids[i], amount, detail, user)
				} else {
					panic(errors.New("new transaction detail requires amount"))
				}
			}
		}
	}
	if len(ids) > 0 {
		validateDetails(tx, ids)
	}
	return ids
}

// UpdateTransactions updates transactions.
func UpdateTransactions(tx *sql.Tx, updates []map[string]interface{}, user string) []int64 {
	ids := make([]int64, len(updates))
	for i, transaction := range updates {
		values := database.InputObject(transaction)
		ids[i] = values.RequireInt("id")
		version := values.RequireInt("version")
		updateTransaction(tx, ids[i], version, values, user)
		if details, ok := values["details"]; ok {
			detailUpdates := details.([]map[string]interface{})
			updateTxDetails(tx, ids[i], detailUpdates, user)
		}
	}
	if len(ids) > 0 {
		validateDetails(tx, ids)
	}
	return ids
}

// DeleteTransactions deletes transactions.
func DeleteTransactions(tx *sql.Tx, ids []map[string]interface{}) {
	deleteRelatedDetails(tx, ids)
	deleteTransactionDetails(tx, ids)
	deleteTransactions(tx, ids)
}
