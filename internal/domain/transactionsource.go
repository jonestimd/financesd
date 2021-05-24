package domain

import (
	"database/sql"

	"github.com/jonestimd/financesd/internal/database/table"
)

// transactionSource provides transactions for a GraphQL query.
type transactionSource struct {
	accountID          int64
	txIDs              []int64
	detailsByTxID      map[int64][]*TransactionDetail
	relatedDetailsByID map[int64]*TransactionDetail
	relatedTxByID      map[int64]*Transaction
}

func (ts *transactionSource) setSource(dbTransactions []*table.Transaction) []*Transaction {
	transactions := make([]*Transaction, len(dbTransactions))
	for i, tx := range dbTransactions {
		transactions[i] = &Transaction{source: ts, Transaction: tx}
	}
	return transactions
}

// loadTransactionDetails loads transaction details grouped by transaction ID.
func (ts *transactionSource) loadTransactionDetails(tx *sql.Tx) {
	if ts.detailsByTxID == nil {
		if ts.txIDs != nil {
			ts.setDetails(getDetailsByTxIDs(tx, ts.txIDs))
		} else {
			ts.setDetails(getDetailsByAccountID(tx, ts.accountID))
		}
	}
}

func (ts *transactionSource) setDetails(details []*table.TransactionDetail) {
	ts.detailsByTxID = make(map[int64][]*TransactionDetail)
	for _, dbDetail := range details {
		if _, ok := ts.detailsByTxID[dbDetail.TransactionID]; !ok {
			ts.detailsByTxID[dbDetail.TransactionID] = make([]*TransactionDetail, 0, 1)
		}
		detail := &TransactionDetail{txSource: ts, TransactionDetail: dbDetail}
		ts.detailsByTxID[dbDetail.TransactionID] = append(ts.detailsByTxID[dbDetail.TransactionID], detail)
	}
}

// loadRelatedDetails loads all related details for the query.
func (ts *transactionSource) loadRelatedDetails(tx *sql.Tx) {
	if ts.relatedDetailsByID == nil {
		if ts.txIDs != nil {
			ts.setRelatedDetails(getRelatedDetailsByTxIDs(tx, ts.txIDs))
		} else {
			ts.setRelatedDetails(getRelatedDetailsByAccountID(tx, ts.accountID))
		}
	}
}

func (ts *transactionSource) setRelatedDetails(details []*table.TransactionDetail) {
	ts.relatedDetailsByID = make(map[int64]*TransactionDetail, len(details))
	for _, detail := range details {
		ts.relatedDetailsByID[detail.ID] = &TransactionDetail{TransactionDetail: detail, txSource: ts}
	}
}

// loadRelatedTransactions loads all related transactions for the account.
func (ts *transactionSource) loadRelatedTransactions(tx *sql.Tx) {
	if ts.relatedTxByID == nil {
		if ts.txIDs != nil {
			ts.setRelatedTransactions(getRelatedTransactions(tx, ts.txIDs))
		} else {
			ts.setRelatedTransactions(getRelatedTransactionsByAccountID(tx, ts.accountID))
		}
	}
}

func (ts *transactionSource) setRelatedTransactions(dbTransactions []*table.Transaction) {
	transactions := ts.setSource(dbTransactions)
	ts.relatedTxByID = make(map[int64]*Transaction, len(transactions))
	for _, tx := range transactions {
		ts.relatedTxByID[tx.ID] = tx
	}
}
