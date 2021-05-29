package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/jonestimd/financesd/internal/database/table"
)

var detailType = reflect.TypeOf(table.TransactionDetail{})

func runDetailQuery(tx *sql.Tx, query string, args ...interface{}) []*table.TransactionDetail {
	rows := runQuery(tx, detailType, query, args...)
	return rows.([]*table.TransactionDetail)
}

const accountTxDetailsSQL = `select td.*
from transaction t
join transaction_detail td on t.id = td.transaction_id
where t.account_id = ?
order by t.id, td.id`

func GetDetailsByAccountID(tx *sql.Tx, accountID int64) []*table.TransactionDetail {
	return runDetailQuery(tx, accountTxDetailsSQL, accountID)
}

const txDetailsSQL = `select td.*
from transaction_detail td
where json_contains(?, cast(transaction_id as json))
order by td.transaction_id, td.id`

func GetDetailsByTxIDs(tx *sql.Tx, txIDs []int64) []*table.TransactionDetail {
	return runDetailQuery(tx, txDetailsSQL, int64sToJson(txIDs))
}

const accountRelatedDetailsSQL = `select rd.*
from transaction tx
join transaction_detail td on tx.id = td.transaction_id
join transaction_detail rd on td.related_detail_id = rd.id
where tx.account_id = ?`

func GetRelatedDetailsByAccountID(tx *sql.Tx, accountID int64) []*table.TransactionDetail {
	return runDetailQuery(tx, accountRelatedDetailsSQL, accountID)
}

const relatedDetailsSQL = `select rd.*
from transaction_detail td
join transaction_detail rd on td.related_detail_id = rd.id
where json_contains(?, cast(td.transaction_id as json))`

func GetRelatedDetailsByTxIDs(tx *sql.Tx, txIDs []int64) []*table.TransactionDetail {
	return runDetailQuery(tx, relatedDetailsSQL, int64sToJson(txIDs))
}

const insertDetailSQL = `insert into transaction_detail
(transaction_id, amount, transaction_category_id, transaction_group_id, memo, asset_quantity, exchange_asset_id, change_date, change_user, version)
values (?, ?, ?, ?, ?, ?, current_timestamp, ?, 0)`

func InsertDetail(tx *sql.Tx, txID int64, amount interface{}, values InputObject, user string) {
	id := runInsert(tx, insertDetailSQL, txID, amount, values.IntOrNull("transactionCategoryId"), values.IntOrNull("transactionGroupId"),
		values.StringOrNull("memo"), values.FloatOrNull("assetQuantity"), values.IntOrNull("exchangeAssetId"), user)
	if transferAccountId, setTransfer := values.GetInt("transferAccountId"); setTransfer {
		insertTransferDetail(tx, id, transferAccountId, user)
	}
}

const insertTransferTransactionSQL = `insert into transaction (account_id, date, payee_id, cleared, change_date, change_user, version)
select ?, rt.date, rt.payee_id, 'N', current_timestamp, ?, 0
from transaction_detail rd
join transaction rt on rd.transaction_id = rt.id
where rd.id = ?`

const insertTransferDetailSQL = `insert into transaction_detail (transaction_id, related_detail_id, amount, change_date, change_user, version)
select ?, ?, -amount, current_timestamp, ?, 0
from transaction_detail
where id = ?`

const setRelatedDetailSQL = `update transaction_detail set related_detail_id = ? where id = ?`

var insertTransferDetail = func(tx *sql.Tx, relatedDetailId int64, accountId interface{}, user string) {
	txId := runInsert(tx, insertTransferTransactionSQL, accountId, user, relatedDetailId)
	detailId := runInsert(tx, insertTransferDetailSQL, txId, relatedDetailId, user, relatedDetailId)
	runUpdate(tx, setRelatedDetailSQL, detailId, relatedDetailId)
}

const moveTransferDetailSQL = `update transaction set account_id = ?, change_date = current_timestamp, change_user = ?, version = version+1
where id = (select d.transaction_id from transaction_detail d where d.related_detail_id = ?)`

var AddOrUpdateTransfer = func(tx *sql.Tx, relatedDetailId int64, accountId interface{}, user string) {
	count := runUpdate(tx, moveTransferDetailSQL, accountId, user, relatedDetailId)
	if count == 0 {
		insertTransferDetail(tx, relatedDetailId, accountId, user)
	}
}

const setTransferAmountSQL = `update transaction_detail set amount = -?, change_date = current_timestamp, change_user = ?, version = version+1
where related_detail_id = ?`

func SetTransferAmount(tx *sql.Tx, relatedDetailId int64, amount interface{}, user string) {
	runUpdate(tx, setTransferAmountSQL, amount, user, relatedDetailId)
}

const updateTxDetailSQL = `update transaction_detail
set amount = case when ? then ? else amount end
, transaction_category_id = case when ? then ? else transaction_category_id end
, transaction_group_id = case when ? then ? else transaction_group_id end
, memo = case when ? then ? else memo end
, asset_quantity = case when ? then ? else asset_quantity end
, exchange_asset_id = case when ? then ? else exchange_asset_id end
, change_date = current_timestamp, change_user = ?, version = version+1
where id = ? and version = ?`

func UpdateDetail(tx *sql.Tx, id int64, version int64, setCategory bool, categoryId interface{}, values InputObject, user string) {
	amount, setAmount := values.GetFloat("amount")
	groupId, setGroup := values.GetInt("transactionGroupId")
	memo, setMemo := values.GetString("memo")
	shares, setShares := values.GetFloat("assetQuantity")
	securityId, setSecurity := values.GetInt("exchangeAssetId")
	count := runUpdate(tx, updateTxDetailSQL,
		setAmount, amount,
		setCategory, categoryId,
		setGroup, groupId,
		setMemo, memo,
		setShares, shares,
		setSecurity, securityId,
		user, id, version)
	if count == 0 {
		panic(fmt.Errorf("transaction detail not found (%d @ %d)", id, version))
	}
}

const deleteEmptyTransactionsSQL = `delete from transaction
where not exists(select 1 from transaction_detail where transaction_id = transaction.id)`

const deleteDetailsSQL = `delete from transaction_detail
where json_contains(?, json_object('ID', id, 'Version', version))`

func DeleteDetails(tx *sql.Tx, ids []*VersionID) {
	deleteIDs, _ := json.Marshal(ids)
	if count := runUpdate(tx, deleteDetailsSQL, deleteIDs); int(count) != len(ids) {
		panic(errors.New("transaction detail(s) not found"))
	}
	runUpdate(tx, deleteEmptyTransactionsSQL)
}

const deleteRelatedDetailSQL = `delete from transaction_detail
where related_detail_id in (select id from transaction_detail where json_contains(?, json_object('id', transaction_id)))`

func DeleteRelatedDetails(tx *sql.Tx, txIDs []map[string]interface{}) {
	deleteIDs, _ := json.Marshal(txIDs)
	if count := runUpdate(tx, deleteRelatedDetailSQL, deleteIDs); count > 0 {
		runUpdate(tx, deleteEmptyTransactionsSQL)
	}
}

func DeleteTransactionDetails(tx *sql.Tx, txIDs []map[string]interface{}) {
	deleteIDs, _ := json.Marshal(txIDs)
	runUpdate(tx, "delete from transaction_detail where json_contains(?, json_object('id', transaction_id))", deleteIDs)
}

const deleteTransferDetailSQL = `delete from transaction_detail
where id = (select related_detail_id from transaction_detail where id = ?)`

func DeleteTransfer(tx *sql.Tx, relatedDetailId int64) {
	if count := runUpdate(tx, deleteTransferDetailSQL, relatedDetailId); count != 0 {
		runUpdate(tx, deleteEmptyTransactionsSQL)
	}
}

const validateDetailsSQL = `select id, error
from (
    select td.id
    , case when t.security_id is null and tc.security = 'Y' is null then concat('security required for category: ', tc.id)
        when tc.asset_exchange = 'Y' and td.asset_quantity is null then concat('shares required for category: ', tc.id)
        when tc.asset_exchange = 'N' and td.asset_quantity is not null then concat('shares not allowed for category: ', tc.id)
        when tc.income = 'N' and td.asset_quantity < 0 then 'shares must be positive for expense category'
        when tc.income = 'Y' and td.asset_quantity > 0 then 'shares must be negative for income category'
        else null
	  end error
    from transaction_detail td
    join transaction t on td.transaction_id = t.id
    left join transaction_category tc on td.transaction_category_id = tc.id
	where json_contains(?, td.transaction_id)) errors
where errors.error is not null`

// ValidateDetails checks for invalid security fields and returns a map of detail ID to validation error.
func ValidateDetails(tx *sql.Tx, transactionIDs []int64) {
	rows, err := tx.Query(validateDetailsSQL, int64sToJson(transactionIDs))
	if err != nil {
		panic(err)
	}
	result := make(map[int64]string)
	for rows.Next() {
		var id int64
		var text string
		rows.Scan(&id, &text)
		result[id] = text
	}
	if len(result) > 0 {
		panic(fmt.Errorf("transaction detail errors: %v", result))
	}
}
