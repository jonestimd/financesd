package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// TransactionDetail represents a line item of a financial transaction.
type TransactionDetail struct {
	ID                    int64
	TransactionID         int64
	Transaction           *Transaction
	TransactionCategoryID *int64
	TransactionGroupID    *int64
	Memo                  *string
	Amount                float64
	AssetQuantity         *float64
	ExchangeAssetID       *int64
	RelatedDetailID       *int64
	RelatedDetail         *TransactionDetail
	Version               int
	txSource              *transactionSource
	Audited
}

func (d *TransactionDetail) ptrTo(column string) interface{} {
	switch column {
	case "id":
		return &d.ID
	case "transaction_id":
		return &d.TransactionID
	case "transaction_category_id":
		return &d.TransactionCategoryID
	case "transaction_group_id":
		return &d.TransactionGroupID
	case "memo":
		return &d.Memo
	case "amount":
		return &d.Amount
	case "asset_quantity":
		return &d.AssetQuantity
	case "exchange_asset_id":
		return &d.ExchangeAssetID
	case "related_detail_id":
		return &d.RelatedDetailID
	case "version":
		return &d.Version
	}
	return d.Audited.ptrToAudit(column)
}

var detailType = reflect.TypeOf(TransactionDetail{})

// GetRelatedDetail returns the related detail.
func (d *TransactionDetail) GetRelatedDetail(tx *sql.Tx) (*TransactionDetail, error) {
	if d.RelatedDetailID == nil {
		return nil, nil
	}
	if d.txSource.loadRelatedDetails(tx) != nil {
		return nil, d.txSource.err
	}
	return d.txSource.relatedDetailsByID[*d.RelatedDetailID], nil
}

// GetRelatedTransaction returns the transaction for a related detail.
func (d *TransactionDetail) GetRelatedTransaction(tx *sql.Tx) (*Transaction, error) {
	if d.txSource.loadRelatedTransactions(tx) != nil {
		return nil, d.txSource.err
	}
	return d.txSource.relatedTxByID[d.TransactionID], nil
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

var updateTxDetails = func(tx *sql.Tx, txId int64, details []map[string]interface{}, user string) error {
	deleteIDs := make([]*versionID, 0)
	for _, detail := range details {
		values := inputObject(detail)
		if id, err := values.getVersionID(); err != nil {
			return err
		} else if id != nil {
			if len(detail) == 2 {
				deleteIDs = append(deleteIDs, id)
			} else {
				amount, setAmount := values.getFloat("amount")
				transferAccountId, setTransfer := values.getInt("transferAccountId")
				categoryId, setCategory := values.getInt("transactionCategoryId")
				if setTransfer && setCategory {
					return errors.New("cannot specify both transferAccountId and transactionCategoryId")
				}
				if setTransfer {
					setCategory = true
					categoryId = nil
				}
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
					user, id.ID, id.Version)
				if err != nil {
					return err
				} else if count == 0 {
					return fmt.Errorf("transaction detail not found (%d @ %d)", id.ID, id.Version)
				}
				if setCategory && transferAccountId == nil {
					if err := deleteTransfer(tx, id.ID); err != nil {
						return err
					}
				} else if setTransfer {
					if err := addOrUpdateTransfer(tx, id.ID, transferAccountId, user); err != nil {
						return err
					}
				}
				if setAmount && !setCategory {
					return setTransferAmount(tx, id.ID, amount, user)
				}
			}
		} else {
			if amount, ok := values.getFloat("amount"); ok && amount != nil {
				if err := insertDetail(tx, amount, values, user); err != nil {
					return err
				}
			} else {
				return errors.New("amount is required to add a transaction detail")
			}
		}
	}
	if len(deleteIDs) > 0 {
		return deleteDetails(tx, deleteIDs)
	}
	return nil
}

const deleteDetailsSQL = `delete from transaction_detail
where json_contains(?, json_object('ID', id, 'Version', version))`

func deleteDetails(tx *sql.Tx, ids []*versionID) error {
	deleteIDs, _ := json.Marshal(ids)
	if count, err := runUpdate(tx, deleteDetailsSQL, deleteIDs); err != nil {
		return err
	} else if int(count) != len(ids) {
		return errors.New("transaction detail(s) not found")
	}
	_, err := runUpdate(tx, deleteEmptyTransactionsSQL)
	return err
}

const insertDetailSQL = `insert into transaction_detail
(amount, transaction_category_id, transaction_group_id, memo, asset_quantity, exchange_asset_id, change_date, change_user, version)
values (?, ?, ?, ?, ?, current_timestamp, ?, 0)`

func insertDetail(tx *sql.Tx, amount interface{}, values inputObject, user string) error {
	id, err := runInsert(tx, insertDetailSQL, amount, values.intOrNull("transactionCategoryId"), values.intOrNull("transactionGroupId"),
		values.stringOrNull("memo"), values.floatOrNull("assetQuantity"), values.intOrNull("exchangeAssetId"), user)
	if err != nil {
		return err
	}
	if transferAccountId, setTransfer := values.getInt("transferAccountId"); setTransfer {
		return insertTransferDetail(tx, id, transferAccountId, user)
	}
	return nil
}

const deleteTransferDetailSQL = `delete from transaction_detail
where id = (select related_detail_id from transaction_detail where id = ?)`

const deleteEmptyTransactionsSQL = `delete from transaction
where not exists(select 1 from transaction_detail where transaction_id = transaction.id)`

var deleteTransfer = func(tx *sql.Tx, relatedDetailId int64) error {
	if count, err := runUpdate(tx, deleteTransferDetailSQL, relatedDetailId); err != nil {
		return err
	} else if count == 0 {
		return nil
	}
	_, err := runUpdate(tx, deleteEmptyTransactionsSQL)
	return err
}

const moveTransferDetailSQL = `update transaction set account_id = ?, change_date = current_timestamp, change_user = ?, version = version+1
where id = (select d.transaction_id from transaction_detail d where d.related_detail_id = ?)`

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

func insertTransferDetail(tx *sql.Tx, relatedDetailId int64, accountId interface{}, user string) error {
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

var addOrUpdateTransfer = func(tx *sql.Tx, relatedDetailId int64, accountId interface{}, user string) error {
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

func setTransferAmount(tx *sql.Tx, relatedDetailId int64, amount interface{}, user string) error {
	_, err := runUpdate(tx, setTransferAmountSQL, amount, user, relatedDetailId)
	return err
}
