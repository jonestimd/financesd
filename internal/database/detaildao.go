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

func runDetailQuery(tx *sql.Tx, query string, args ...interface{}) ([]*table.TransactionDetail, error) {
	rows, err := runQuery(tx, detailType, query, args...)
	if err != nil {
		return nil, err
	}
	return rows.([]*table.TransactionDetail), nil
}

const accountTxDetailsSQL = `select td.*
from transaction t
join transaction_detail td on t.id = td.transaction_id
where t.account_id = ?
order by t.id, td.id`

func GetDetailsByAccountID(tx *sql.Tx, accountID int64) ([]*table.TransactionDetail, error) {
	return runDetailQuery(tx, accountTxDetailsSQL, accountID)
}

const txDetailsSQL = `select td.*
from transaction_detail td
where json_contains(?, cast(transaction_id as json))
order by t.id, td.id`

func GetDetailsByTxIDs(tx *sql.Tx, txIDs []int64) ([]*table.TransactionDetail, error) {
	return runDetailQuery(tx, txDetailsSQL, int64sToJson(txIDs))
}

const accountRelatedDetailsSQL = `select rd.*
from transaction tx
join transaction_detail td on tx.id = td.transaction_id
join transaction_detail rd on td.related_detail_id = rd.id
where tx.account_id = ?`

func GetRelatedDetailsByAccountID(tx *sql.Tx, accountID int64) ([]*table.TransactionDetail, error) {
	return runDetailQuery(tx, accountRelatedDetailsSQL, accountID)
}

const relatedDetailsSQL = `select rd.*
from transaction_detail td
join transaction_detail rd on td.related_detail_id = rd.id
where json_contains(?, cast(td.transaction_id as json))`

func GetRelatedDetailsByTxIDs(tx *sql.Tx, txIDs []int64) ([]*table.TransactionDetail, error) {
	return runDetailQuery(tx, relatedDetailsSQL, int64sToJson(txIDs))
}

const insertDetailSQL = `insert into transaction_detail
(amount, transaction_category_id, transaction_group_id, memo, asset_quantity, exchange_asset_id, change_date, change_user, version)
values (?, ?, ?, ?, ?, current_timestamp, ?, 0)`

func InsertDetail(tx *sql.Tx, amount interface{}, values InputObject, user string) error {
	id, err := runInsert(tx, insertDetailSQL, amount, values.IntOrNull("transactionCategoryId"), values.IntOrNull("transactionGroupId"),
		values.StringOrNull("memo"), values.FloatOrNull("assetQuantity"), values.IntOrNull("exchangeAssetId"), user)
	if err != nil {
		return err
	}
	if transferAccountId, setTransfer := values.GetInt("transferAccountId"); setTransfer {
		return insertTransferDetail(tx, id, transferAccountId, user)
	}
	return nil
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

var insertTransferDetail = func(tx *sql.Tx, relatedDetailId int64, accountId interface{}, user string) error {
	if txId, err := runInsert(tx, insertTransferTransactionSQL, accountId, user, relatedDetailId); err != nil {
		return err
	} else {
		if detailId, err := runInsert(tx, insertTransferDetailSQL, txId, relatedDetailId, user, relatedDetailId); err != nil {
			return err
		} else {
			_, err := runUpdate(tx, setRelatedDetailSQL, detailId, relatedDetailId)
			return err
		}
	}
}

const moveTransferDetailSQL = `update transaction set account_id = ?, change_date = current_timestamp, change_user = ?, version = version+1
where id = (select d.transaction_id from transaction_detail d where d.related_detail_id = ?)`

var AddOrUpdateTransfer = func(tx *sql.Tx, relatedDetailId int64, accountId interface{}, user string) error {
	count, err := runUpdate(tx, moveTransferDetailSQL, accountId, user, relatedDetailId)
	if err != nil {
		return err
	}
	if count == 0 {
		return insertTransferDetail(tx, relatedDetailId, accountId, user)
	}
	return nil
}

const setTransferAmountSQL = `update transaction_detail set amount = -?, change_date = current_timestamp, change_user = ?, version = version+1
where related_detail_id = ?`

func SetTransferAmount(tx *sql.Tx, relatedDetailId int64, amount interface{}, user string) error {
	_, err := runUpdate(tx, setTransferAmountSQL, amount, user, relatedDetailId)
	return err
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

func UpdateDetail(tx *sql.Tx, id int64, version int64, setCategory bool, categoryId interface{}, values InputObject, user string) error {
	amount, setAmount := values.GetFloat("amount")
	groupId, setGroup := values.GetInt("transactionGroupId")
	memo, setMemo := values.GetString("memo")
	shares, setShares := values.GetFloat("assetQuantity")
	securityId, setSecurity := values.GetInt("exchangeAssetId")
	count, err := runUpdate(tx, updateTxDetailSQL,
		setAmount, amount,
		setCategory, categoryId,
		setGroup, groupId,
		setMemo, memo,
		setShares, shares,
		setSecurity, securityId,
		user, id, version)
	if err != nil {
		return err
	} else if count == 0 {
		return fmt.Errorf("transaction detail not found (%d @ %d)", id, version)
	}
	return nil
}

const deleteEmptyTransactionsSQL = `delete from transaction
where not exists(select 1 from transaction_detail where transaction_id = transaction.id)`

const deleteDetailsSQL = `delete from transaction_detail
where json_contains(?, json_object('ID', id, 'Version', version))`

func DeleteDetails(tx *sql.Tx, ids []*VersionID) error {
	deleteIDs, _ := json.Marshal(ids)
	if count, err := runUpdate(tx, deleteDetailsSQL, deleteIDs); err != nil {
		return err
	} else if int(count) != len(ids) {
		return errors.New("transaction detail(s) not found")
	}
	_, err := runUpdate(tx, deleteEmptyTransactionsSQL)
	return err
}

const deleteTransferDetailSQL = `delete from transaction_detail
where id = (select related_detail_id from transaction_detail where id = ?)`

func DeleteTransfer(tx *sql.Tx, relatedDetailId int64) error {
	if count, err := runUpdate(tx, deleteTransferDetailSQL, relatedDetailId); err != nil {
		return err
	} else if count == 0 {
		return nil
	}
	_, err := runUpdate(tx, deleteEmptyTransactionsSQL)
	return err
}
