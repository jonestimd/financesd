package database

import (
	"database/sql"
	"fmt"
)

type inputObject interface {
	getInt(key string) (interface{}, bool)
	getFloat(key string) (interface{}, bool)
	getString(key string) (interface{}, bool)
}

const accountRelatedDetailsSQL = `select rd.*
	from transaction tx
	join transaction_detail td on tx.id = td.transaction_id
	join transaction_detail rd on td.related_detail_id = rd.id
	where tx.account_id = ?`

func GetRelatedDetailsByAccountID(tx *sql.Tx, accountID int64) ([]*TransactionDetail, error) {
	if rows, err := runQuery(tx, detailType, accountRelatedDetailsSQL, accountID); err != nil {
		return nil, err
	} else {
		details := rows.([]*TransactionDetail)
		return details, err
	}
}

const txDetailsSQL = `select td.*
	from transaction_detail td
	where json_contains(?, cast(transaction_id as json))
	order by t.id, td.id`

func GetRelatedDetailsByTxIDs(tx *sql.Tx, txIDs []int64) ([]*TransactionDetail, error) {
	if rows, err := runQuery(tx, detailType, txDetailsSQL, int64sToJson(txIDs)); err != nil {
		return nil, err
	} else {
		details := rows.([]*TransactionDetail)
		return details, err
	}
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

func UpdateDetail(tx *sql.Tx, id int64, version int64, values inputObject, user string) error {
	amount, setAmount := values.getFloat("amount")
	categoryId, setCategory := values.getInt("transactionCategoryId")
	groupId, setGroup := values.getInt("transactionGroupId")
	memo, setMemo := values.getString("memo")
	shares, setShares := values.getFloat("assetQuantity")
	securityId, setSecurity := values.getInt("exchangeAssetId")
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
