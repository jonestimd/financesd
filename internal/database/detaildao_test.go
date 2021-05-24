package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MonsantoCo/mocka/v2"
	"github.com/jonestimd/financesd/internal/database/table"
	"github.com/jonestimd/financesd/internal/sqltest"
	"github.com/stretchr/testify/assert"
)

func testDetailsQuery(t *testing.T, doQuery func(tx *sql.Tx) ([]*table.TransactionDetail, string, []interface{})) {
	details := []*table.TransactionDetail{{ID: 96}}
	sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
		runQueryStub := mocka.Function(t, &runQuery, details)
		defer runQueryStub.Restore()

		result, query, params := doQuery(tx)

		assert.Equal(t, details, result)
		assert.Equal(t, []interface{}{tx, detailType, query, params}, runQueryStub.GetCall(0).Arguments())
	})
}

func Test_GetDetailsByAccountID(t *testing.T) {
	testDetailsQuery(t, func(tx *sql.Tx) ([]*table.TransactionDetail, string, []interface{}) {
		accountID := int64(42)

		result := GetDetailsByAccountID(tx, accountID)

		return result, accountTxDetailsSQL, []interface{}{accountID}
	})
}

func Test_GetDetailsByTxIDs(t *testing.T) {
	testDetailsQuery(t, func(tx *sql.Tx) ([]*table.TransactionDetail, string, []interface{}) {
		txIDs := []int64{42}

		result := GetDetailsByTxIDs(tx, txIDs)

		return result, txDetailsSQL, []interface{}{int64sToJson(txIDs)}
	})
}

func Test_GetRelatedDetailsByAccountID(t *testing.T) {
	testDetailsQuery(t, func(tx *sql.Tx) ([]*table.TransactionDetail, string, []interface{}) {
		accountID := int64(42)

		result := GetRelatedDetailsByAccountID(tx, accountID)

		return result, accountRelatedDetailsSQL, []interface{}{accountID}
	})
}

func Test_GetRelatedDetailsByTxIDs(t *testing.T) {
	testDetailsQuery(t, func(tx *sql.Tx) ([]*table.TransactionDetail, string, []interface{}) {
		txIDs := []int64{42}

		result := GetRelatedDetailsByTxIDs(tx, txIDs)

		return result, relatedDetailsSQL, []interface{}{int64sToJson(txIDs)}
	})
}

func Test_InsertDetail(t *testing.T) {
	id := int64(42)
	user := "user id"
	amount := 420.0
	t.Run("inserts detail with category", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			values := InputObject{
				"transactionCategoryId": 96,
				"transactionGroupId":    69,
				"memo":                  "notes",
				"assetQuantity":         4.2,
				"exchangeAssetId":       24,
			}
			runInsertStub := mocka.Function(t, &runInsert, id)
			defer runInsertStub.Restore()

			InsertDetail(tx, amount, values, user)

			assert.Equal(t, sqltest.UpdateArgs(tx, insertDetailSQL, amount, int64(96), int64(69), "notes", 4.2, int64(24), user),
				runInsertStub.GetCall(0).Arguments())
		})
	})
	t.Run("inserts transfer", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			values := InputObject{"transferAccountId": 96}
			runInsertStub := mocka.Function(t, &runInsert, id)
			defer runInsertStub.Restore()
			insertTransferStub := mocka.Function(t, &insertTransferDetail)
			defer insertTransferStub.Restore()

			InsertDetail(tx, amount, values, user)

			assert.Equal(t, sqltest.UpdateArgs(tx, insertDetailSQL, amount, nil, nil, nil, nil, nil, user), runInsertStub.GetCall(0).Arguments())
			assert.Equal(t, []interface{}{tx, id, int64(96), user}, insertTransferStub.GetCall(0).Arguments())
		})
	})
}

func Test_insertTransferDetail(t *testing.T) {
	txID := int64(42)
	detailID := int64(96)
	relatedDetailID := int64(69)
	accountID := int64(123)
	user := "user id"
	sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
		runInsertStub := mocka.Function(t, &runInsert, txID)
		runInsertStub.OnCall(1).Return(detailID)
		defer runInsertStub.Restore()
		runUpdateStub := mocka.Function(t, &runUpdate, int64(1))
		defer runUpdateStub.Restore()

		insertTransferDetail(tx, relatedDetailID, accountID, user)

		assert.Equal(t, sqltest.UpdateArgs(tx, insertTransferTransactionSQL, accountID, user, relatedDetailID), runInsertStub.GetCall(0).Arguments())
		assert.Equal(t, sqltest.UpdateArgs(tx, insertTransferDetailSQL, txID, relatedDetailID, user, relatedDetailID), runInsertStub.GetCall(1).Arguments())
		assert.Equal(t, sqltest.UpdateArgs(tx, setRelatedDetailSQL, detailID, relatedDetailID), runUpdateStub.GetCall(0).Arguments())
	})
}

func Test_AddOrUpdateTransfer(t *testing.T) {
	relatedDetailID := int64(69)
	accountID := int64(123)
	user := "user id"
	t.Run("updates existing transfer", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			runUpdateStub := mocka.Function(t, &runUpdate, int64(1))
			defer runUpdateStub.Restore()

			AddOrUpdateTransfer(tx, relatedDetailID, accountID, user)

			assert.Equal(t, sqltest.UpdateArgs(tx, moveTransferDetailSQL, accountID, user, relatedDetailID), runUpdateStub.GetCall(0).Arguments())
		})
	})
	t.Run("inserts transfer", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			runUpdateStub := mocka.Function(t, &runUpdate, int64(0))
			defer runUpdateStub.Restore()
			insertTransferDetailStub := mocka.Function(t, &insertTransferDetail)
			defer insertTransferDetailStub.Restore()

			AddOrUpdateTransfer(tx, relatedDetailID, accountID, user)

			assert.Equal(t, sqltest.UpdateArgs(tx, moveTransferDetailSQL, accountID, user, relatedDetailID), runUpdateStub.GetCall(0).Arguments())
			assert.Equal(t, []interface{}{tx, relatedDetailID, accountID, user}, insertTransferDetailStub.GetCall(0).Arguments())
		})
	})
}

func Test_SetTransferAmount_returnsError(t *testing.T) {
	relatedDetailID := int64(69)
	amount := 420.0
	user := "user id"
	sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
		runUpdateStub := mocka.Function(t, &runUpdate, int64(0))
		defer runUpdateStub.Restore()

		SetTransferAmount(tx, relatedDetailID, amount, user)

		assert.Equal(t, sqltest.UpdateArgs(tx, setTransferAmountSQL, amount, user, relatedDetailID), runUpdateStub.GetCall(0).Arguments())
	})
}

func Test_UpdateDetail(t *testing.T) {
	id := int64(42)
	version := int64(1)
	categoryID := int64(96)
	user := "user id"
	t.Run("updates the detail", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			values := InputObject{
				"amount":             42.0,
				"transactionGroupId": 69,
				"memo":               "notes",
				"assetQuantity":      4.2,
				"exchangeAssetId":    24,
			}
			runUpdateStub := mocka.Function(t, &runUpdate, int64(1))
			defer runUpdateStub.Restore()

			UpdateDetail(tx, id, version, true, categoryID, values, user)

			assert.Equal(t, sqltest.UpdateArgs(
				tx, updateTxDetailSQL, true, 42.0, true, categoryID, true, int64(69), true, "notes", true, 4.2, true, int64(24), user, id, version),
				runUpdateStub.GetCall(0).Arguments())
		})
	})
	t.Run("panics for detail not found", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			values := InputObject{"memo": "notes"}
			runUpdateStub := mocka.Function(t, &runUpdate, int64(0))
			defer runUpdateStub.Restore()
			defer func() {
				assert.Equal(t, sqltest.UpdateArgs(
					tx, updateTxDetailSQL, false, nil, false, nil, false, nil, true, "notes", false, nil, false, nil, user, id, version),
					runUpdateStub.GetCall(0).Arguments())
				if err := recover(); err != nil {
					assert.Equal(t, "transaction detail not found (42 @ 1)", err.(error).Error())
				} else {
					assert.Fail(t, "expected an error")
				}
			}()

			UpdateDetail(tx, id, version, false, nil, values, user)
		})
	})
}

func Test_DeleteDetails(t *testing.T) {
	ids := []*VersionID{{42, 1}, {24, 0}}
	t.Run("delete details and empty transactions", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			runUpdateStub := mocka.Function(t, &runUpdate, int64(2))
			runUpdateStub.OnCall(1).Return(int64(0))
			defer runUpdateStub.Restore()

			DeleteDetails(tx, ids)

			assert.Equal(t, 2, runUpdateStub.CallCount())
			idArg, _ := json.Marshal(ids)
			assert.Equal(t, sqltest.UpdateArgs(tx, deleteDetailsSQL, idArg), runUpdateStub.GetCall(0).Arguments())
			assert.Equal(t, sqltest.UpdateArgs(tx, deleteEmptyTransactionsSQL), runUpdateStub.GetCall(1).Arguments())
		})
	})
	t.Run("panics for detail not found", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			runUpdateStub := mocka.Function(t, &runUpdate, int64(0))
			defer runUpdateStub.Restore()
			defer func() {
				assert.Equal(t, 1, runUpdateStub.CallCount())
				if err := recover(); err != nil {
					assert.Equal(t, "transaction detail(s) not found", err.(error).Error())
				} else {
					assert.Fail(t, "expected an error")
				}
			}()

			DeleteDetails(tx, ids)
		})
	})
}

func Test_DeleteTransfer(t *testing.T) {
	relatedDetailID := int64(42)
	t.Run("delete transfer detail and empty transactions", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			runUpdateStub := mocka.Function(t, &runUpdate, int64(2))
			runUpdateStub.OnCall(1).Return(int64(0))
			defer runUpdateStub.Restore()

			DeleteTransfer(tx, relatedDetailID)

			assert.Equal(t, 2, runUpdateStub.CallCount())
			assert.Equal(t, sqltest.UpdateArgs(tx, deleteTransferDetailSQL, relatedDetailID), runUpdateStub.GetCall(0).Arguments())
			assert.Equal(t, sqltest.UpdateArgs(tx, deleteEmptyTransactionsSQL), runUpdateStub.GetCall(1).Arguments())
		})
	})
	t.Run("ignores transfer detail not found", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			runUpdateStub := mocka.Function(t, &runUpdate, int64(0))
			defer runUpdateStub.Restore()

			DeleteTransfer(tx, relatedDetailID)

			assert.Equal(t, 1, runUpdateStub.CallCount())
		})
	})
}

func Test_ValidateDetails(t *testing.T) {
	transactionIDs := []int64{42, 96, 69}
	t.Run("runs validation query", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			rows := sqltest.MockRows("id", "error").AddRow(int64(42), "invalid shares").AddRow(int64(96), "shares required")
			mockDB.ExpectQuery(validateDetailsSQL).WithArgs(int64sToJson(transactionIDs)).WillReturnRows(rows)

			result := ValidateDetails(tx, transactionIDs)

			assert.Equal(t, map[int64]string{42: "invalid shares", 96: "shares required"}, result)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})
	})
	t.Run("panics for query error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			expectedErr := errors.New("query error")
			mockDB.ExpectQuery(validateDetailsSQL).WithArgs(int64sToJson(transactionIDs)).WillReturnError(expectedErr)
			defer func() {
				if err := recover(); err != nil {
					assert.Same(t, expectedErr, err)
				} else {
					assert.Fail(t, "expected an error")
				}
			}()

			ValidateDetails(tx, transactionIDs)
		})
	})
}
