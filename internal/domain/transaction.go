package domain

import (
	"database/sql"
	"fmt"

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
	return graphql.DefaultResolveFn(replaceSource(p, t.Transaction))
}

// GetDetails returns details for a transaction in the account.
func (t *Transaction) GetDetails(tx *sql.Tx) ([]*TransactionDetail, error) {
	if t.source.loadTransactionDetails(tx) != nil {
		return nil, t.source.err
	}
	return t.source.detailsByTxID[t.ID], nil
}

// GetTransactions returns all transactions for the account.
func GetTransactions(tx *sql.Tx, accountID int64) ([]*Transaction, error) {
	at := &transactionSource{accountID: accountID}
	return at.setSource(getTransactions(tx, accountID))
}

// GetTransactionsByIDs returns transactions for the specified IDs.
func GetTransactionsByIDs(tx *sql.Tx, ids []int64) ([]*Transaction, error) {
	source := &transactionSource{txIDs: ids}
	return source.setSource(getTransactionsByIDs(tx, ids))
}

// UpdateTransactions updates transactions.
func UpdateTransactions(tx *sql.Tx, updates []map[string]interface{}, user string) ([]int64, error) {
	ids := make([]int64, len(updates))
	for i, transaction := range updates {
		values := database.InputObject(transaction)
		ids[i] = values.RequireInt("id")
		version := values.RequireInt("version")
		if err := updateTransaction(tx, ids[i], version, values, user); err != nil {
			return nil, err
		}
		if details, ok := values["details"]; ok {
			detailUpdates := details.([]map[string]interface{})
			if err := updateTxDetails(tx, ids[i], detailUpdates, user); err != nil {
				return nil, err
			}
		}
	}
	if len(ids) > 0 {
		if errors, err := validateDetails(tx, ids); err != nil {
			return nil, err
		} else if len(errors) > 0 {
			return nil, fmt.Errorf("transaction detail errors: %v", errors)
		}
	}
	return ids, nil
}
