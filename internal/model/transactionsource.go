package model

import "database/sql"

// transactionSource provides transactions for a GraphQL query.
type transactionSource struct {
	accountID          int64
	txIDs              []int64
	err                error
	detailsByTxID      map[int64][]*TransactionDetail
	relatedDetailsByID map[int64]*TransactionDetail
	relatedTxByID      map[int64]*Transaction
}

func (at *transactionSource) selectQuery(byAccountID string, byTxIDs string) (string, interface{}) {
	if at.txIDs != nil {
		return byTxIDs, int64sToJson(at.txIDs)
	}
	return byAccountID, at.accountID
}

const accountTxDetailsSQL = `select td.*
	from transaction t
	join transaction_detail td on t.id = td.transaction_id
	where t.account_id = ?
	order by t.id, td.id`

const txDetailsSQL = `select td.*
	from transaction_detail td
	where json_contains(?, cast(transaction_id as json))
	order by t.id, td.id`

// loadTransactionDetails loads transaction details grouped by transaction ID.
func (at *transactionSource) loadTransactionDetails(tx *sql.Tx) error {
	if at.err == nil && at.detailsByTxID == nil {
		query, param := at.selectQuery(accountTxDetailsSQL, txDetailsSQL)
		var rows interface{}
		if rows, at.err = runQuery(tx, detailType, query, param); at.err == nil {
			details := rows.([]*TransactionDetail)
			at.detailsByTxID = make(map[int64][]*TransactionDetail)
			for _, detail := range details {
				detail.txSource = at
				if _, ok := at.detailsByTxID[detail.TransactionID]; !ok {
					at.detailsByTxID[detail.TransactionID] = make([]*TransactionDetail, 0, 1)
				}
				at.detailsByTxID[detail.TransactionID] = append(at.detailsByTxID[detail.TransactionID], detail)
			}
		}
	}
	return at.err
}

const accountRelatedDetailsSQL = `select rd.*
	from transaction tx
	join transaction_detail td on tx.id = td.transaction_id
	join transaction_detail rd on td.related_detail_id = rd.id
	where tx.account_id = ?`

const relatedDetailsSQL = `select rd.*
	from transaction_detail td
	join transaction_detail rd on td.related_detail_id = rd.id
	where json_contains(?, cast(td.transaction_id as json))`

// loadRelatedDetails loads all related details for the query.
func (at *transactionSource) loadRelatedDetails(tx *sql.Tx) error {
	if at.err == nil && at.relatedDetailsByID == nil {
		query, param := at.selectQuery(accountRelatedDetailsSQL, relatedDetailsSQL)
		var rows interface{}
		if rows, at.err = runQuery(tx, detailType, query, param); at.err == nil {
			details := rows.([]*TransactionDetail)
			at.relatedDetailsByID = make(map[int64]*TransactionDetail, len(details))
			for _, detail := range details {
				detail.txSource = at
				at.relatedDetailsByID[detail.ID] = detail
			}
		}
	}
	return at.err
}

const accountRelatedTxSQL = `select rt.*
	from transaction tx
	join transaction_detail td on tx.id = td.transaction_id
	join transaction_detail rd on td.related_detail_id = rd.id
	join transaction rt on rd.transaction_id = rt.id
	where tx.account_id = ?`

const relatedTxSQL = `select rt.*
	from transaction_detail td
	join transaction_detail rd on td.related_detail_id = rd.id
	join transaction rt on rd.transaction_id = rt.id
	where json_contains(?, cast(td.transaction_id as json))`

// loadRelatedTransactions loads all related transactions for the account.
func (at *transactionSource) loadRelatedTransactions(tx *sql.Tx) error {
	if at.err == nil && at.relatedTxByID == nil {
		query, param := at.selectQuery(accountRelatedTxSQL, relatedTxSQL)
		var rows interface{}
		if rows, at.err = runQuery(tx, transactionType, query, param); at.err == nil {
			transactions := rows.([]*Transaction)
			at.relatedTxByID = make(map[int64]*Transaction, len(transactions))
			for _, tx := range transactions {
				at.relatedTxByID[tx.ID] = tx
			}
		}
	}
	return at.err
}
