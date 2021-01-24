package model

import (
	"database/sql"
	"reflect"
	"time"
)

// Transaction represents a financial transaction.
type Transaction struct {
	ID                  int64
	Date                time.Time
	Memo                *string
	ReferenceNumber     *string
	Cleared             *YesNo
	AccountID           int64
	PayeeID             *int64
	SecurityID          *int64
	Details             []TransactionDetail
	Version             int
	accountTransactions *AccountTransactions
	Audited
}

func (t *Transaction) ptrTo(column string) interface{} {
	switch column {
	case "id":
		return &t.ID
	case "date":
		return &t.Date
	case "memo":
		return &t.Memo
	case "reference_number":
		return &t.ReferenceNumber
	case "cleared":
		return &t.Cleared
	case "account_id":
		return &t.AccountID
	case "payee_id":
		return &t.PayeeID
	case "security_id":
		return &t.SecurityID
	case "version":
		return &t.Version
	}
	return t.Audited.ptrToAudit(column)
}

// AccountTransactions contains transactions for an account.
type AccountTransactions struct {
	accountID          int64
	err                error
	detailsByTxID      map[int64][]*TransactionDetail
	relatedDetailsByID map[int64]*TransactionDetail
	relatedTxByID      map[int64]*Transaction
}

var transactionType = reflect.TypeOf(Transaction{})

const transactionSQL = "select * from transaction where account_id = ? order by date, id"

// GetTransactions returns all transactions for the account.
func GetTransactions(tx *sql.Tx, accountID int64) ([]*Transaction, error) {
	at := &AccountTransactions{accountID: accountID}
	rows, err := runQuery(tx, transactionType, transactionSQL, accountID)
	if err != nil {
		return nil, err
	}
	transactions := rows.([]*Transaction)
	for _, transaction := range transactions {
		transaction.accountTransactions = at
	}
	return transactions, nil
}

const transactionDetailSQL = `select td.*
	from transaction t
	join transaction_detail td on t.id = td.transaction_id
	where t.account_id = ?
	order by t.id, td.id`

// getTransactionDetails loads transaction details grouped by transaction ID.
func (at *AccountTransactions) getTransactionDetails(tx *sql.Tx) error {
	if rows, err := runQuery(tx, detailType, transactionDetailSQL, at.accountID); err != nil {
		at.err = err
	} else {
		details := rows.([]*TransactionDetail)
		at.detailsByTxID = make(map[int64][]*TransactionDetail)
		for _, detail := range details {
			detail.accountTransactions = at
			if _, ok := at.detailsByTxID[detail.TransactionID]; !ok {
				at.detailsByTxID[detail.TransactionID] = make([]*TransactionDetail, 0, 1)
			}
			at.detailsByTxID[detail.TransactionID] = append(at.detailsByTxID[detail.TransactionID], detail)
		}
	}
	return at.err
}

// GetDetails returns details for a transaction in the account.
func (t *Transaction) GetDetails(tx *sql.Tx) ([]*TransactionDetail, error) {
	if t.accountTransactions.err != nil {
		return nil, t.accountTransactions.err
	}
	if t.accountTransactions.detailsByTxID == nil && t.accountTransactions.getTransactionDetails(tx) != nil {
		return nil, t.accountTransactions.err
	}
	return t.accountTransactions.detailsByTxID[t.ID], nil
}

const relatedDetailsSQL = `select rd.*
	from transaction tx
	join transaction_detail td on tx.id = td.transaction_id
	join transaction_detail rd on td.related_detail_id = rd.id
	where tx.account_id = ?`

// getRelatedDetails loads all related details for the account.
func (at *AccountTransactions) getRelatedDetails(tx *sql.Tx) error {
	if rows, err := runQuery(tx, detailType, relatedDetailsSQL, at.accountID); err != nil {
		at.err = err
	} else {
		details := rows.([]*TransactionDetail)
		at.relatedDetailsByID = make(map[int64]*TransactionDetail, len(details))
		for _, detail := range details {
			detail.accountTransactions = at
			at.relatedDetailsByID[detail.ID] = detail
		}
	}
	return at.err
}

const relatedTransactionSQL = `select rt.*
	from transaction tx
	join transaction_detail td on tx.id = td.transaction_id
	join transaction_detail rd on td.related_detail_id = rd.id
	join transaction rt on rd.transaction_id = rt.id
	where tx.account_id = ?`

// getRelatedTransactions loads all related transactions for the account.
func (at *AccountTransactions) getRelatedTransactions(tx *sql.Tx) error {
	if rows, err := runQuery(tx, transactionType, relatedTransactionSQL, at.accountID); err != nil {
		at.err = err
	} else {
		transactions := rows.([]*Transaction)
		at.relatedTxByID = make(map[int64]*Transaction, len(transactions))
		for _, tx := range transactions {
			at.relatedTxByID[tx.ID] = tx
		}
	}
	return at.err
}
