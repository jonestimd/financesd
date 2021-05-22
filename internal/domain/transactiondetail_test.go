package domain

import (
	"database/sql"
	"errors"
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
	expectedErr := errors.New("test error")
	relatedDetail := NewTransactionDetail(id, -1)
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
			detail := &TransactionDetail{TransactionDetail: &table.TransactionDetail{RelatedDetailID: test.relatedID}, txSource: test.txSource}

			result, err := detail.GetRelatedDetail(nil)

			assert.Equal(t, test.expectedErr, err)
			assert.Equal(t, test.expectedResult, result)
		})
	}
}

func Test_GetRelatedTransaction(t *testing.T) {
	id := int64(42)
	expectedErr := errors.New("test error")
	relatedTx := NewTransaction(id)
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
			detail := &TransactionDetail{TransactionDetail: &table.TransactionDetail{TransactionID: id}, txSource: test.txSource}

			result, err := detail.GetRelatedTransaction(nil)

			assert.Equal(t, test.expectedErr, err)
			assert.Equal(t, test.expectedResult, result)
		})
	}
}

func Test_updateTxDetails_update(t *testing.T) {
	user := "user id"
	id := 42
	txID := int64(69)
	version := 1
	t.Run("returns update detail error", func(t *testing.T) {
		update := database.InputObject{
			"id":                    id,
			"version":               version,
			"transactionCategoryId": -1,
		}
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			expectedErr := errors.New("update failed")
			updateDetailStub := mocka.Function(t, &updateDetail, expectedErr)
			defer updateDetailStub.Restore()

			err := updateTxDetails(tx, txID, []map[string]interface{}{update}, user)

			assert.Same(t, expectedErr, err)
		})
	})
	t.Run("returns error for id without version", func(t *testing.T) {
		update := database.InputObject{
			"id":     id,
			"amount": 42.0,
		}
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			err := updateTxDetails(tx, txID, []map[string]interface{}{update}, user)

			assert.Equal(t, "version is required for update/delete", err.Error())
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
				updateDetailStub := mocka.Function(t, &updateDetail, nil)
				defer updateDetailStub.Restore()
				deleteTransferStub := mocka.Function(t, &deleteTransfer, nil)
				defer deleteTransferStub.Restore()
				setTransferAmountStub := mocka.Function(t, &setTransferAmount, nil)
				defer setTransferAmountStub.Restore()

				err := updateTxDetails(tx, txID, []map[string]interface{}{update}, user)

				assert.Nil(t, err)
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
				updateDetailStub := mocka.Function(t, &updateDetail, nil)
				defer updateDetailStub.Restore()
				deleteTransferStub := mocka.Function(t, &deleteTransfer, nil)
				defer deleteTransferStub.Restore()
				setTransferAmountStub := mocka.Function(t, &setTransferAmount, nil)
				defer setTransferAmountStub.Restore()

				err := updateTxDetails(tx, txID, []map[string]interface{}{update}, user)

				assert.Nil(t, err)
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
		t.Run("deletes transfer when category is set", func(t *testing.T) {
			sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
				updateDetailStub := mocka.Function(t, &updateDetail, nil)
				defer updateDetailStub.Restore()
				deleteTransferStub := mocka.Function(t, &deleteTransfer, nil)
				defer deleteTransferStub.Restore()

				err := updateTxDetails(tx, txID, []map[string]interface{}{update}, user)

				assert.Nil(t, err)
				assert.Equal(t, []interface{}{tx, int64(id), int64(version), true, int64(96), update, user}, updateDetailStub.GetCall(0).Arguments())
				assert.Equal(t, []interface{}{tx, int64(id)}, deleteTransferStub.GetCall(0).Arguments())
			})
		})
		t.Run("returns delete transfer error", func(t *testing.T) {
			sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
				expectedErr := errors.New("delete failed")
				updateDetailStub := mocka.Function(t, &updateDetail, nil)
				defer updateDetailStub.Restore()
				deleteTransferStub := mocka.Function(t, &deleteTransfer, expectedErr)
				defer deleteTransferStub.Restore()

				err := updateTxDetails(tx, txID, []map[string]interface{}{update}, user)

				assert.Same(t, expectedErr, err)
			})
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
				updateDetailStub := mocka.Function(t, &updateDetail, nil)
				defer updateDetailStub.Restore()
				deleteTransferStub := mocka.Function(t, &deleteTransfer, nil)
				defer deleteTransferStub.Restore()

				err := updateTxDetails(tx, txID, []map[string]interface{}{update}, user)

				assert.Nil(t, err)
				assert.Equal(t, []interface{}{tx, int64(id)}, deleteTransferStub.GetCall(0).Arguments())
			})
		})
		t.Run("updates transfer", func(t *testing.T) {
			update["transferAccountId"] = 96
			sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
				updateDetailStub := mocka.Function(t, &updateDetail, nil)
				defer updateDetailStub.Restore()
				addOrUpdateTransferStub := mocka.Function(t, &addOrUpdateTransfer, nil)
				defer addOrUpdateTransferStub.Restore()

				err := updateTxDetails(tx, txID, []map[string]interface{}{update}, user)

				assert.Nil(t, err)
				assert.Equal(t, []interface{}{tx, int64(id), int64(96), user}, addOrUpdateTransferStub.GetCall(0).Arguments())
			})
		})
		t.Run("returns update transfer error", func(t *testing.T) {
			update["transferAccountId"] = 96
			sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
				expectedErr := errors.New("update transfer failed")
				updateDetailStub := mocka.Function(t, &updateDetail, nil)
				defer updateDetailStub.Restore()
				addOrUpdateTransferStub := mocka.Function(t, &addOrUpdateTransfer, expectedErr)
				defer addOrUpdateTransferStub.Restore()

				err := updateTxDetails(tx, txID, []map[string]interface{}{update}, user)

				assert.Same(t, expectedErr, err)
			})
		})
	})
	t.Run("returns error for transfer with category", func(t *testing.T) {
		update := database.InputObject{
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
		update := database.InputObject{
			"id":      id,
			"version": version,
			"amount":  42.0,
		}
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			updateDetailStub := mocka.Function(t, &updateDetail, expectedErr)
			defer updateDetailStub.Restore()

			err := updateTxDetails(tx, txID, []map[string]interface{}{update}, user)

			assert.Equal(t, expectedErr, err)
		})
	})
}

func Test_updateTxDetails_insert(t *testing.T) {
	user := "user id"
	txID := int64(69)
	t.Run("requires amount", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			update := database.InputObject{}

			err := updateTxDetails(tx, txID, []map[string]interface{}{update}, user)

			assert.Equal(t, "amount is required to add a transaction detail", err.Error())
		})
	})
	t.Run("inserts a detail", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			update := database.InputObject{"amount": 42.0}
			insertDetailStub := mocka.Function(t, &insertDetail, nil)
			defer insertDetailStub.Restore()

			err := updateTxDetails(tx, txID, []map[string]interface{}{update}, user)

			assert.Nil(t, err)
			assert.Equal(t, []interface{}{tx, 42.0, update, user}, insertDetailStub.GetCall(0).Arguments())
		})
	})
	t.Run("returns insert error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			update := database.InputObject{"amount": 42.0}
			expectedErr := errors.New("insert failed")
			insertDetailStub := mocka.Function(t, &insertDetail, expectedErr)
			defer insertDetailStub.Restore()

			err := updateTxDetails(tx, txID, []map[string]interface{}{update}, user)

			assert.Same(t, expectedErr, err)
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
	versionID, _ := update.GetVersionID()
	t.Run("deletes details", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			deleteDetailsStub := mocka.Function(t, &deleteDetails, nil)
			defer deleteDetailsStub.Restore()

			err := updateTxDetails(tx, txID, []map[string]interface{}{update}, user)

			assert.Nil(t, err)
			assert.Equal(t, []interface{}{tx, []*database.VersionID{versionID}}, deleteDetailsStub.GetCall(0).Arguments())
		})
	})
	t.Run("returns deletes error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			expectedErr := errors.New("delete failed")
			deleteDetailsStub := mocka.Function(t, &deleteDetails, expectedErr)
			defer deleteDetailsStub.Restore()

			err := updateTxDetails(tx, txID, []map[string]interface{}{update}, user)

			assert.Same(t, expectedErr, err)
		})
	})
}
