package domain

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MonsantoCo/mocka/v2"
	"github.com/graphql-go/graphql"
	"github.com/jonestimd/financesd/internal/database"
	"github.com/jonestimd/financesd/internal/database/table"
	"github.com/jonestimd/financesd/internal/sqltest"
	"github.com/stretchr/testify/assert"
)

func Test_TransactionDetail_Resolve(t *testing.T) {
	expectedResult := "result"
	defaultResolveStub := mocka.Function(t, &defaultResolveFn, expectedResult, nil)
	defer defaultResolveStub.Restore()
	p := graphql.ResolveParams{}
	detail := NewTransactionDetail(42, 1)

	result, err := detail.Resolve(p)

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
	assert.Same(t, detail.TransactionDetail, defaultResolveStub.GetCall(0).Arguments()[0].(graphql.ResolveParams).Source)
}

func Test_GetRelatedDetail(t *testing.T) {
	id := int64(42)
	relatedDetail := NewTransactionDetail(id, -1)
	relatedDetailsByID := map[int64]*TransactionDetail{id: relatedDetail}
	tests := []struct {
		name           string
		relatedID      *int64
		txSource       *transactionSource
		expectedResult *TransactionDetail
	}{
		{"returns nil for nil RelatedDetailID", nil, nil, nil},
		{"returns existing detail", &id, &transactionSource{relatedDetailsByID: relatedDetailsByID}, relatedDetail},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			detail := &TransactionDetail{TransactionDetail: &table.TransactionDetail{RelatedDetailID: test.relatedID}, txSource: test.txSource}

			result := detail.GetRelatedDetail(nil)

			assert.Equal(t, test.expectedResult, result)
		})
	}
}

func Test_GetRelatedTransaction(t *testing.T) {
	id := int64(42)
	relatedTx := NewTransaction(id)
	relatedTxByID := map[int64]*Transaction{id: relatedTx}
	txSource := &transactionSource{relatedTxByID: relatedTxByID}
	detail := &TransactionDetail{TransactionDetail: &table.TransactionDetail{TransactionID: id}, txSource: txSource}

	result := detail.GetRelatedTransaction(nil)

	assert.Equal(t, relatedTx, result)
}

func Test_updateTxDetails_update(t *testing.T) {
	user := "user id"
	id := 42
	txID := int64(69)
	version := 1
	t.Run("panics for id without version", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			update := database.InputObject{
				"id":     id,
				"amount": 42.0,
			}
			defer func() {
				if err := recover(); err != nil {
					assert.Equal(t, "version is required for update/delete", err.(error).Error())
				} else {
					assert.Fail(t, "expected an error")
				}
			}()

			updateTxDetails(tx, txID, []map[string]interface{}{update}, user)
		})
	})
	t.Run("update amount", func(t *testing.T) {
		t.Run("does not update transfer if category is set", func(t *testing.T) {
			update := database.InputObject{
				"id":                    id,
				"version":               version,
				"transactionCategoryId": nil,
				"amount":                42.0,
			}
			sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
				updateDetailStub := mocka.Function(t, &updateDetail)
				defer updateDetailStub.Restore()
				deleteTransferStub := mocka.Function(t, &deleteTransfer)
				defer deleteTransferStub.Restore()
				setTransferAmountStub := mocka.Function(t, &setTransferAmount)
				defer setTransferAmountStub.Restore()

				updateTxDetails(tx, txID, []map[string]interface{}{update}, user)

				assert.Equal(t, []interface{}{tx, int64(id), int64(version), true, nil, update, user}, updateDetailStub.GetCall(0).Arguments())
				assert.Equal(t, []interface{}{tx, int64(id)}, deleteTransferStub.GetCall(0).Arguments())
				assert.Equal(t, 0, setTransferAmountStub.CallCount())
			})
		})
		t.Run("updates transfer if category not set", func(t *testing.T) {
			update := database.InputObject{
				"id":      id,
				"version": version,
				"amount":  42.0,
			}
			sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
				updateDetailStub := mocka.Function(t, &updateDetail)
				defer updateDetailStub.Restore()
				deleteTransferStub := mocka.Function(t, &deleteTransfer)
				defer deleteTransferStub.Restore()
				setTransferAmountStub := mocka.Function(t, &setTransferAmount)
				defer setTransferAmountStub.Restore()

				updateTxDetails(tx, txID, []map[string]interface{}{update}, user)

				assert.Equal(t, []interface{}{tx, int64(id), int64(version), false, nil, update, user}, updateDetailStub.GetCall(0).Arguments())
				assert.Equal(t, 0, deleteTransferStub.CallCount())
				assert.Equal(t, []interface{}{tx, int64(id), 42.0, user}, setTransferAmountStub.GetCall(0).Arguments())
			})
		})
	})
	t.Run("update category", func(t *testing.T) {
		update := database.InputObject{
			"id":                    id,
			"version":               version,
			"transactionCategoryId": 96,
		}
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			updateDetailStub := mocka.Function(t, &updateDetail)
			defer updateDetailStub.Restore()
			deleteTransferStub := mocka.Function(t, &deleteTransfer)
			defer deleteTransferStub.Restore()

			updateTxDetails(tx, txID, []map[string]interface{}{update}, user)

			assert.Equal(t, []interface{}{tx, int64(id), int64(version), true, int64(96), update, user}, updateDetailStub.GetCall(0).Arguments())
			assert.Equal(t, []interface{}{tx, int64(id)}, deleteTransferStub.GetCall(0).Arguments())
		})
	})
	t.Run("update transfer account", func(t *testing.T) {
		update := database.InputObject{
			"id":      id,
			"version": version,
		}
		t.Run("deletes transfer", func(t *testing.T) {
			update["transferAccountId"] = nil
			sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
				updateDetailStub := mocka.Function(t, &updateDetail)
				defer updateDetailStub.Restore()
				deleteTransferStub := mocka.Function(t, &deleteTransfer)
				defer deleteTransferStub.Restore()

				updateTxDetails(tx, txID, []map[string]interface{}{update}, user)

				assert.Equal(t, []interface{}{tx, int64(id)}, deleteTransferStub.GetCall(0).Arguments())
			})
		})
		t.Run("updates transfer", func(t *testing.T) {
			update["transferAccountId"] = 96
			sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
				updateDetailStub := mocka.Function(t, &updateDetail)
				defer updateDetailStub.Restore()
				addOrUpdateTransferStub := mocka.Function(t, &addOrUpdateTransfer)
				defer addOrUpdateTransferStub.Restore()

				updateTxDetails(tx, txID, []map[string]interface{}{update}, user)

				assert.Equal(t, []interface{}{tx, int64(id), int64(96), user}, addOrUpdateTransferStub.GetCall(0).Arguments())
			})
		})
	})
	t.Run("panics for transfer with category", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			update := database.InputObject{
				"id":                    id,
				"version":               version,
				"transactionCategoryId": 96,
				"transferAccountId":     69,
			}
			defer func() {
				if err := recover(); err != nil {
					assert.Equal(t, "cannot specify both transferAccountId and transactionCategoryId", err.(error).Error())
				} else {
					assert.Fail(t, "expected an error")
				}
			}()

			updateTxDetails(tx, txID, []map[string]interface{}{update}, user)
		})
	})
}

func Test_updateTxDetails_insert(t *testing.T) {
	user := "user id"
	txID := int64(69)
	t.Run("panics for no amount", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			update := database.InputObject{}
			defer func() {
				if err := recover(); err != nil {
					assert.Equal(t, "new transaction detail requires amount", err.(error).Error())
				} else {
					assert.Fail(t, "expected an error")
				}
			}()

			updateTxDetails(tx, txID, []map[string]interface{}{update}, user)
		})
	})
	t.Run("inserts a detail", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			update := database.InputObject{"amount": 42.0}
			insertDetailStub := mocka.Function(t, &insertDetail)
			defer insertDetailStub.Restore()

			updateTxDetails(tx, txID, []map[string]interface{}{update}, user)

			assert.Equal(t, []interface{}{tx, txID, 42.0, update, user}, insertDetailStub.GetCall(0).Arguments())
		})
	})
}

func Test_updateTxDetails_delete(t *testing.T) {
	user := "user id"
	id := 42
	version := 1
	txID := int64(69)
	update := database.InputObject{
		"id":      id,
		"version": version,
	}
	versionID := update.GetVersionID()
	sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
		deleteDetailsStub := mocka.Function(t, &deleteDetails)
		defer deleteDetailsStub.Restore()

		updateTxDetails(tx, txID, []map[string]interface{}{update}, user)

		assert.Equal(t, []interface{}{tx, []*database.VersionID{versionID}}, deleteDetailsStub.GetCall(0).Arguments())
	})
}
