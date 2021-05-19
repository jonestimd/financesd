package database

import (
	"database/sql"
	"fmt"
)

const accountTransactionsSQL = "select * from transaction where account_id = ? order by date, id"

func runTransactionQuery(tx *sql.Tx, query string, args ...interface{}) ([]*Transaction, error) {
	rows, err := runQuery(tx, transactionType, query, args...)
	if err != nil {
		return nil, err
	}
	return rows.([]*Transaction), nil
}

// GetTransactions returns all transactions for the account.
func GetTransactions(tx *sql.Tx, accountID int64) ([]*Transaction, error) {
	return runTransactionQuery(tx, accountTransactionsSQL, accountID)
}

const transactionsByIDSQL = "select * from transaction where json_contains(?, cast(id as json))"

// GetTransactionsByIDs returns transactions for the specified IDs.
func GetTransactionsByIDs(tx *sql.Tx, ids []int64) ([]*Transaction, error) {
	return runTransactionQuery(tx, transactionsByIDSQL, int64sToJson(ids))
}

const accountRelatedTxSQL = `select rt.*
	from transaction tx
	join transaction_detail td on tx.id = td.transaction_id
	join transaction_detail rd on td.related_detail_id = rd.id
	join transaction rt on rd.transaction_id = rt.id
	where tx.account_id = ?`

// GetRelatedTransactionsByAccountID returns all related transactions for the account.
func GetRelatedTransactionsByAccountID(tx *sql.Tx, accountID int64) ([]*Transaction, error) {
	return runTransactionQuery(tx, accountRelatedTxSQL, accountID)
}

const relatedTxSQL = `select rt.*
	from transaction_detail td
	join transaction_detail rd on td.related_detail_id = rd.id
	join transaction rt on rd.transaction_id = rt.id
	where json_contains(?, cast(td.transaction_id as json))`

// GetRelatedTransactions returns all related transactions for the transaction IDs.
func GetRelatedTransactions(tx *sql.Tx, relatedTxIDs []int64) ([]*Transaction, error) {
	return runTransactionQuery(tx, relatedTxSQL, int64sToJson(relatedTxIDs))
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
func UpdateTransaction(tx *sql.Tx, id int64, version int64, values InputObject, user string) error {
	ref, setRef := values.GetString("referenceNumber")
	payeeId, setPayee := values.GetInt("payeeId")
	securityId, setSecurity := values.GetInt("securityId")
	memo, setMemo := values.GetString("memo")
	accountId, setAccount := values.GetInt("accountId")
	count, err := runUpdate(tx, updateTxSQL,
		values.StringOrNull("date"),
		setRef, ref,
		setPayee, payeeId,
		setSecurity, securityId,
		setMemo, memo,
		values.StringOrNull("cleared"),
		setAccount, accountId,
		user, id, version)
	if err != nil {
		return err
	} else if count == 0 {
		return fmt.Errorf("transaction not found (%d @ %d)", id, version)
	}
	return nil
}
