package model

import (
	"database/sql"
	"fmt"
	"reflect"
	"time"
)

// Transaction represents a financial transaction.
type Transaction struct {
	ID              int64
	Date            time.Time
	Memo            *string
	ReferenceNumber *string
	Cleared         *YesNo
	AccountID       int64
	PayeeID         *int64
	SecurityID      *int64
	Details         []TransactionDetail
	Version         int
	source          *transactionSource
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

var transactionType = reflect.TypeOf(Transaction{})

const accountTransactionsSQL = "select * from transaction where account_id = ? order by date, id"

// GetTransactions returns all transactions for the account.
func GetTransactions(tx *sql.Tx, accountID int64) ([]*Transaction, error) {
	at := &transactionSource{accountID: accountID}
	rows, err := runQuery(tx, transactionType, accountTransactionsSQL, accountID)
	if err != nil {
		return nil, err
	}
	transactions := rows.([]*Transaction)
	for _, transaction := range transactions {
		transaction.source = at
	}
	return transactions, nil
}

const transactionsByIDSQL = "select * from transaction where json_contains(?, cast(id as json))"

// GetDetails returns details for a transaction in the account.
func (t *Transaction) GetDetails(tx *sql.Tx) ([]*TransactionDetail, error) {
	if t.source.loadTransactionDetails(tx) != nil {
		return nil, t.source.err
	}
	return t.source.detailsByTxID[t.ID], nil
}

const updateTxSQL = `update transaction
set date = coalesce(?, date)
, reference_number = case when ? then ? else reference_number end
, payee_id = case when ? then ? else payee_id end
, security_id = case when ? then ? else security_id end
, memo = case when ? then ? else memo end
, cleared = coalesce(?, cleared)
, account_id = case when ? then ? else account_id end
, change_date = current_timestamp, change_user = ?, version = version+1
where id = ? and version = ?`

// UpdateTransactions updates transactions.
func UpdateTransactions(tx *sql.Tx, args interface{}, user string) ([]int64, error) {
	updates := args.([]interface{})
	ids := make([]int64, len(updates))
	for i, transaction := range updates {
		values := transaction.(map[string]interface{})
		ids[i] = int64(values["id"].(int))
		ref, setRef := values["referenceNumber"]
		payeeId, setPayee := values["payeeId"]
		securityId, setSecurity := values["securityId"]
		memo, setMemo := values["memo"]
		accountId, setAccount := values["accountId"]
		version := int64(values["version"].(int))
		rs, err := runUpdate(tx, updateTxSQL,
			stringOrNull(values["date"]),
			setRef, stringOrNull(ref),
			setPayee, int64OrNull(payeeId),
			setSecurity, int64OrNull(securityId),
			setMemo, stringOrNull(memo),
			stringOrNull(values["cleared"]),
			setAccount, int64OrNull(accountId),
			user, ids[i], version)
		if err != nil {
			return nil, err
		}
		if count, err := rs.RowsAffected(); err != nil {
			return nil, err
		} else if count == 0 {
			return nil, fmt.Errorf("transaction not found (%d) or incorrect version (%d)", ids[i], version)
		}
	}
	return ids, nil
}

// GetTransactionsByIDs returns transactions for the specified IDs.
func GetTransactionsByIDs(tx *sql.Tx, ids []int64) ([]*Transaction, error) {
	rows, err := runQuery(tx, transactionType, transactionsByIDSQL, int64sToJson(ids))
	if err != nil {
		return nil, err
	}
	transactions := rows.([]*Transaction)
	source := &transactionSource{txIDs: ids}
	for _, transaction := range transactions {
		transaction.source = source
	}
	return transactions, nil
}
