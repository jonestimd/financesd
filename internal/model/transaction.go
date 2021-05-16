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
func UpdateTransactions(tx *sql.Tx, updates []map[string]interface{}, user string) ([]int64, error) {
	ids := make([]int64, len(updates))
	for i, transaction := range updates {
		values := inputObject(transaction)
		ids[i] = values.requireInt("id")
		ref, setRef := values.getString("referenceNumber")
		payeeId, setPayee := values.getInt("payeeId")
		securityId, setSecurity := values.getInt("securityId")
		memo, setMemo := values.getString("memo")
		accountId, setAccount := values.getInt("accountId")
		version := int64(values["version"].(int))
		count, err := runUpdate(tx, updateTxSQL,
			values.stringOrNull("date"),
			setRef, ref,
			setPayee, payeeId,
			setSecurity, securityId,
			setMemo, memo,
			values.stringOrNull("cleared"),
			setAccount, accountId,
			user, ids[i], version)
		if err != nil {
			return nil, err
		} else if count == 0 {
			return nil, fmt.Errorf("transaction not found (%d @ %d)", ids[i], version)
		}
		if details, ok := values["details"]; ok {
			detailUpdates := details.([]map[string]interface{})
			if err := updateTxDetails(tx, ids[i], detailUpdates, user); err != nil {
				return nil, err
			}
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
