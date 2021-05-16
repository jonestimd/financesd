package model

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MonsantoCo/mocka/v2"
	"github.com/jonestimd/financesd/internal/sqltest"
	"github.com/stretchr/testify/assert"
)

func Test_TransactionDetail_ptrTo(t *testing.T) {
	detail := &TransactionDetail{}
	tests := []struct {
		column string
		ptr    interface{}
	}{
		{column: "id", ptr: &detail.ID},
		{column: "transaction_id", ptr: &detail.TransactionID},
		{column: "transaction_category_id", ptr: &detail.TransactionCategoryID},
		{column: "transaction_group_id", ptr: &detail.TransactionGroupID},
		{column: "memo", ptr: &detail.Memo},
		{column: "amount", ptr: &detail.Amount},
		{column: "asset_quantity", ptr: &detail.AssetQuantity},
		{column: "exchange_asset_id", ptr: &detail.ExchangeAssetID},
		{column: "related_detail_id", ptr: &detail.RelatedDetailID},
		{column: "version", ptr: &detail.Version},
		{column: "change_user", ptr: &detail.ChangeUser},
		{column: "change_date", ptr: &detail.ChangeDate},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s returns pointer to field", test.column), func(t *testing.T) {
			field := detail.ptrTo(test.column)

			assert.Same(t, test.ptr, field)
		})
	}
}

func Test_GetRelatedDetail(t *testing.T) {
	id := int64(42)
	expectedErr := errors.New("test error")
	relatedDetail := &TransactionDetail{ID: id}
	relatedDetailsByID := map[int64]*TransactionDetail{id: relatedDetail}
	tests := []struct {
		name           string
		relatedID      *int64
		txSource       *transactionSource
		expectedErr    error
		expectedResult *TransactionDetail
	}{
		{"returns nil for nil RelatedDetailID", nil, nil, nil, nil},
		{"returns existing error", &id, &transactionSource{err: expectedErr}, expectedErr, nil},
		{"returns existing detail", &id, &transactionSource{relatedDetailsByID: relatedDetailsByID}, nil, relatedDetail},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			detail := &TransactionDetail{RelatedDetailID: test.relatedID, txSource: test.txSource}

			result, err := detail.GetRelatedDetail(nil)

			assert.Equal(t, test.expectedErr, err)
			assert.Equal(t, test.expectedResult, result)
		})
	}
}

func Test_GetRelatedTransaction(t *testing.T) {
	id := int64(42)
	expectedErr := errors.New("test error")
	relatedTx := &Transaction{ID: id}
	relatedTxByID := map[int64]*Transaction{id: relatedTx}
	tests := []struct {
		name           string
		txSource       *transactionSource
		expectedErr    error
		expectedResult *Transaction
	}{
		{"returns existing error", &transactionSource{err: expectedErr}, expectedErr, nil},
		{"returns existing detail", &transactionSource{relatedTxByID: relatedTxByID}, nil, relatedTx},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			detail := &TransactionDetail{TransactionID: id, txSource: test.txSource}

			result, err := detail.GetRelatedTransaction(nil)

			assert.Equal(t, test.expectedErr, err)
			assert.Equal(t, test.expectedResult, result)
		})
	}
}

func Test_updateTxDetails(t *testing.T) {
	user := "user id"
	id := 42
	txID := int64(69)
	version := 1
	fieldTests := []struct {
		name  string
		field string
		value interface{}
	}{
		{"sets amount", "amount", 42.0},
		{"sets category", "transactionCategoryId", 96},
		{"sets group", "transactionGroupId", 96},
		{"sets memo", "memo", "notes"},
		{"sets shares", "assetQuantity", 42.0},
		{"sets security", "exchangeAssetId", 96},
	}
	for _, test := range fieldTests {
		t.Run(test.name, func(t *testing.T) {
			update := inputObject{
				"id":      id,
				"version": version,
			}
			update[test.field] = test.value
			sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
				expectedArgs := []interface{}{
					tx, updateTxDetailSQL,
					[]interface{}{
						test.field == "amount", update["amount"],
						test.field == "transactionCategoryId", update.intOrNull("transactionCategoryId"),
						test.field == "transactionGroupId", update.intOrNull("transactionGroupId"),
						test.field == "memo", update["memo"],
						test.field == "assetQuantity", update.floatOrNull("assetQuantity"),
						test.field == "exchangeAssetId", update.intOrNull("exchangeAssetId"),
						user, int64(id), int64(version)}}
				runUpdateStub := mocka.Function(t, &runUpdate, int64(1), nil)
				defer runUpdateStub.Restore()

				err := updateTxDetails(tx, txID, []map[string]interface{}{update}, user)

				assert.Nil(t, err)
				assert.Equal(t, expectedArgs, runUpdateStub.GetCall(0).Arguments())
				assert.Nil(t, mockDB.ExpectationsWereMet())
			})
		})
	}
	t.Run("returns update error for not found", func(t *testing.T) {
		update := inputObject{
			"id":      id,
			"version": version,
			"amount":  42.0,
		}
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			runUpdateStub := mocka.Function(t, &runUpdate, int64(0), nil)
			defer runUpdateStub.Restore()

			err := updateTxDetails(tx, txID, []map[string]interface{}{update}, user)

			assert.Equal(t, "transaction detail not found (42 @ 1)", err.Error())
		})
	})
	t.Run("returns error for id without version", func(t *testing.T) {
		update := inputObject{
			"id":     id,
			"amount": 42.0,
		}
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			err := updateTxDetails(tx, txID, []map[string]interface{}{update}, user)

			assert.Equal(t, "version is required for update/delete", err.Error())
		})
	})
	t.Run("update category", func(t *testing.T) {
		update := inputObject{
			"id":                    id,
			"version":               version,
			"transactionCategoryId": 96,
		}
		t.Run("deletes transfer when category is set", func(t *testing.T) {
			sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
				runUpdateStub := mocka.Function(t, &runUpdate, int64(1), nil)
				deleteTransferStub := mocka.Function(t, &deleteTransfer, nil)
				defer runUpdateStub.Restore()
				defer deleteTransferStub.Restore()

				err := updateTxDetails(tx, txID, []map[string]interface{}{update}, user)

				assert.Nil(t, err)
				assert.Equal(t, []interface{}{tx, int64(id)}, deleteTransferStub.GetCall(0).Arguments())
			})
		})
		t.Run("returns delete transfer error", func(t *testing.T) {
			sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
				expectedErr := errors.New("delete failed")
				runUpdateStub := mocka.Function(t, &runUpdate, int64(1), nil)
				deleteTransferStub := mocka.Function(t, &deleteTransfer, expectedErr)
				defer runUpdateStub.Restore()
				defer deleteTransferStub.Restore()

				err := updateTxDetails(tx, txID, []map[string]interface{}{update}, user)

				assert.Same(t, expectedErr, err)
			})
		})
	})
	t.Run("update transfer account", func(t *testing.T) {
		update := inputObject{
			"id":      id,
			"version": version,
		}
		t.Run("deletes transfer", func(t *testing.T) {
			update["transferAccountId"] = nil
			sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
				runUpdateStub := mocka.Function(t, &runUpdate, int64(1), nil)
				deleteTransferStub := mocka.Function(t, &deleteTransfer, nil)
				defer runUpdateStub.Restore()
				defer deleteTransferStub.Restore()

				err := updateTxDetails(tx, txID, []map[string]interface{}{update}, user)

				assert.Nil(t, err)
				assert.Equal(t, []interface{}{tx, int64(id)}, deleteTransferStub.GetCall(0).Arguments())
			})
		})
		t.Run("updates transfer", func(t *testing.T) {
			update["transferAccountId"] = 96
			sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
				runUpdateStub := mocka.Function(t, &runUpdate, int64(1), nil)
				addOrUpdateTransferStub := mocka.Function(t, &addOrUpdateTransfer, nil)
				defer runUpdateStub.Restore()
				defer addOrUpdateTransferStub.Restore()

				err := updateTxDetails(tx, txID, []map[string]interface{}{update}, user)

				assert.Nil(t, err)
				assert.Equal(t, []interface{}{tx, int64(id), int64(96), user}, addOrUpdateTransferStub.GetCall(0).Arguments())
			})
		})
	})
	t.Run("returns error for transfer with category", func(t *testing.T) {
		update := inputObject{
			"id":                    id,
			"version":               version,
			"transactionCategoryId": 96,
			"transferAccountId":     69,
		}
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			err := updateTxDetails(tx, txID, []map[string]interface{}{update}, user)

			assert.Equal(t, "cannot specify both transferAccountId and transactionCategoryId", err.Error())
		})
	})
	t.Run("returns update error", func(t *testing.T) {
		expectedErr := errors.New("update failed")
		update := inputObject{
			"id":      id,
			"version": version,
			"amount":  42.0,
		}
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			runUpdateStub := mocka.Function(t, &runUpdate, int64(0), expectedErr)
			defer runUpdateStub.Restore()

			err := updateTxDetails(tx, txID, []map[string]interface{}{update}, user)

			assert.Equal(t, expectedErr, err)
		})
	})
}
