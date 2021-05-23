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

func testDetailsQuery(t *testing.T, doQuery func(tx *sql.Tx) ([]*table.TransactionDetail, string, []interface{}, error)) {
	expectedErr := errors.New("query failed")
	tests := []struct {
		name    string
		details []*table.TransactionDetail
		err     error
	}{
		{"returns details", []*table.TransactionDetail{{ID: 96}}, nil},
		{"returns query error", nil, expectedErr},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
				runQueryStub := mocka.Function(t, &runQuery, test.details, test.err)
				defer runQueryStub.Restore()

				result, query, params, err := doQuery(tx)

				assert.Equal(t, test.err, err)
				assert.Equal(t, test.details, result)
				assert.Equal(t, []interface{}{tx, detailType, query, params}, runQueryStub.GetCall(0).Arguments())
			})
		})
	}
}

func Test_GetDetailsByAccountID(t *testing.T) {
	testDetailsQuery(t, func(tx *sql.Tx) ([]*table.TransactionDetail, string, []interface{}, error) {
		accountID := int64(42)

		result, err := GetDetailsByAccountID(tx, accountID)

		return result, accountTxDetailsSQL, []interface{}{accountID}, err
	})
}

func Test_GetDetailsByTxIDs(t *testing.T) {
	testDetailsQuery(t, func(tx *sql.Tx) ([]*table.TransactionDetail, string, []interface{}, error) {
		txIDs := []int64{42}

		result, err := GetDetailsByTxIDs(tx, txIDs)

		return result, txDetailsSQL, []interface{}{int64sToJson(txIDs)}, err
	})
}

func Test_GetRelatedDetailsByAccountID(t *testing.T) {
	testDetailsQuery(t, func(tx *sql.Tx) ([]*table.TransactionDetail, string, []interface{}, error) {
		accountID := int64(42)

		result, err := GetRelatedDetailsByAccountID(tx, accountID)

		return result, accountRelatedDetailsSQL, []interface{}{accountID}, err
	})
}

func Test_GetRelatedDetailsByTxIDs(t *testing.T) {
	testDetailsQuery(t, func(tx *sql.Tx) ([]*table.TransactionDetail, string, []interface{}, error) {
		txIDs := []int64{42}

		result, err := GetRelatedDetailsByTxIDs(tx, txIDs)

		return result, relatedDetailsSQL, []interface{}{int64sToJson(txIDs)}, err
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
			runInsertStub := mocka.Function(t, &runInsert, id, nil)
			defer runInsertStub.Restore()

			err := InsertDetail(tx, amount, values, user)

			assert.Nil(t, err)
			assert.Equal(t, sqltest.UpdateArgs(tx, insertDetailSQL, amount, int64(96), int64(69), "notes", 4.2, int64(24), user),
				runInsertStub.GetCall(0).Arguments())
		})
	})
	t.Run("returns insert error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			values := InputObject{}
			expectedErr := errors.New("insert detail failed")
			runInsertStub := mocka.Function(t, &runInsert, int64(0), expectedErr)
			defer runInsertStub.Restore()

			err := InsertDetail(tx, amount, values, user)

			assert.Same(t, expectedErr, err)
		})
	})
	t.Run("inserts transfer", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			values := InputObject{"transferAccountId": 96}
			runInsertStub := mocka.Function(t, &runInsert, id, nil)
			defer runInsertStub.Restore()
			insertTransferStub := mocka.Function(t, &insertTransferDetail, nil)
			defer insertTransferStub.Restore()

			err := InsertDetail(tx, amount, values, user)

			assert.Nil(t, err)
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
	expectedErr := errors.New("query failed")
	t.Run("inserts transaction and detail", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			runInsertStub := mocka.Function(t, &runInsert, txID, nil)
			runInsertStub.OnCall(1).Return(detailID, nil)
			defer runInsertStub.Restore()
			runUpdateStub := mocka.Function(t, &runUpdate, int64(1), nil)
			defer runUpdateStub.Restore()

			err := insertTransferDetail(tx, relatedDetailID, accountID, user)

			assert.Nil(t, err)
			assert.Equal(t, sqltest.UpdateArgs(tx, insertTransferTransactionSQL, accountID, user, relatedDetailID), runInsertStub.GetCall(0).Arguments())
			assert.Equal(t, sqltest.UpdateArgs(tx, insertTransferDetailSQL, txID, relatedDetailID, user, relatedDetailID), runInsertStub.GetCall(1).Arguments())
			assert.Equal(t, sqltest.UpdateArgs(tx, setRelatedDetailSQL, detailID, relatedDetailID), runUpdateStub.GetCall(0).Arguments())
		})
	})
	t.Run("returns insert transaction error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			runInsertStub := mocka.Function(t, &runInsert, txID, expectedErr)
			defer runInsertStub.Restore()

			err := insertTransferDetail(tx, relatedDetailID, accountID, user)

			assert.Same(t, expectedErr, err)
		})
	})
	t.Run("returns insert detail error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			runInsertStub := mocka.Function(t, &runInsert, txID, nil)
			runInsertStub.OnCall(1).Return(detailID, expectedErr)
			defer runInsertStub.Restore()
			runUpdateStub := mocka.Function(t, &runUpdate, int64(1), nil)
			defer runUpdateStub.Restore()

			err := insertTransferDetail(tx, relatedDetailID, accountID, user)

			assert.Same(t, expectedErr, err)
		})
	})
	t.Run("returns update detail error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			runInsertStub := mocka.Function(t, &runInsert, txID, nil)
			runInsertStub.OnCall(1).Return(detailID, nil)
			defer runInsertStub.Restore()
			runUpdateStub := mocka.Function(t, &runUpdate, int64(0), expectedErr)
			defer runUpdateStub.Restore()

			err := insertTransferDetail(tx, relatedDetailID, accountID, user)

			assert.Same(t, expectedErr, err)
		})
	})
}

func Test_AddOrUpdateTransfer(t *testing.T) {
	relatedDetailID := int64(69)
	accountID := int64(123)
	user := "user id"
	expectedErr := errors.New("query failed")
	t.Run("updates existing transfer", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			runUpdateStub := mocka.Function(t, &runUpdate, int64(1), nil)
			defer runUpdateStub.Restore()

			err := AddOrUpdateTransfer(tx, relatedDetailID, accountID, user)

			assert.Nil(t, err)
			assert.Equal(t, sqltest.UpdateArgs(tx, moveTransferDetailSQL, accountID, user, relatedDetailID), runUpdateStub.GetCall(0).Arguments())
		})
	})
	t.Run("returns update transfer error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			runUpdateStub := mocka.Function(t, &runUpdate, int64(0), expectedErr)
			defer runUpdateStub.Restore()

			err := AddOrUpdateTransfer(tx, relatedDetailID, accountID, user)

			assert.Same(t, expectedErr, err)
		})
	})
	t.Run("inserts transfer", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			runUpdateStub := mocka.Function(t, &runUpdate, int64(0), nil)
			defer runUpdateStub.Restore()
			insertTransferDetailStub := mocka.Function(t, &insertTransferDetail, nil)
			defer insertTransferDetailStub.Restore()

			err := AddOrUpdateTransfer(tx, relatedDetailID, accountID, user)

			assert.Nil(t, err)
			assert.Equal(t, sqltest.UpdateArgs(tx, moveTransferDetailSQL, accountID, user, relatedDetailID), runUpdateStub.GetCall(0).Arguments())
			assert.Equal(t, []interface{}{tx, relatedDetailID, accountID, user}, insertTransferDetailStub.GetCall(0).Arguments())
		})
	})
	t.Run("returns inserts transfer error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			runUpdateStub := mocka.Function(t, &runUpdate, int64(0), nil)
			defer runUpdateStub.Restore()
			insertTransferDetailStub := mocka.Function(t, &insertTransferDetail, expectedErr)
			defer insertTransferDetailStub.Restore()

			err := AddOrUpdateTransfer(tx, relatedDetailID, accountID, user)

			assert.Same(t, expectedErr, err)
		})
	})
}

func Test_SetTransferAmount_returnsError(t *testing.T) {
	relatedDetailID := int64(69)
	amount := 420.0
	user := "user id"
	expectedErr := errors.New("query failed")
	sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
		runUpdateStub := mocka.Function(t, &runUpdate, int64(0), expectedErr)
		defer runUpdateStub.Restore()

		err := SetTransferAmount(tx, relatedDetailID, amount, user)

		assert.Same(t, expectedErr, err)
		assert.Equal(t, sqltest.UpdateArgs(tx, setTransferAmountSQL, amount, user, relatedDetailID), runUpdateStub.GetCall(0).Arguments())
	})
}

func Test_UpdateDetail(t *testing.T) {
	id := int64(42)
	version := int64(1)
	categoryID := int64(96)
	user := "user id"
	expectedErr := errors.New("query failed")
	t.Run("updates the detail", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			values := InputObject{
				"amount":             42.0,
				"transactionGroupId": 69,
				"memo":               "notes",
				"assetQuantity":      4.2,
				"exchangeAssetId":    24,
			}
			runUpdateStub := mocka.Function(t, &runUpdate, int64(1), nil)
			defer runUpdateStub.Restore()

			err := UpdateDetail(tx, id, version, true, categoryID, values, user)

			assert.Nil(t, err)
			assert.Equal(t, sqltest.UpdateArgs(
				tx, updateTxDetailSQL, true, 42.0, true, categoryID, true, int64(69), true, "notes", true, 4.2, true, int64(24), user, id, version),
				runUpdateStub.GetCall(0).Arguments())
		})
	})
	t.Run("returns update error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			values := InputObject{
				"amount": 42.0,
			}
			runUpdateStub := mocka.Function(t, &runUpdate, int64(0), expectedErr)
			defer runUpdateStub.Restore()

			err := UpdateDetail(tx, id, version, false, nil, values, user)

			assert.Same(t, expectedErr, err)
			assert.Equal(t, sqltest.UpdateArgs(
				tx, updateTxDetailSQL, true, 42.0, false, nil, false, nil, false, nil, false, nil, false, nil, user, id, version),
				runUpdateStub.GetCall(0).Arguments())
		})
	})
	t.Run("returns not found error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			values := InputObject{
				"memo": "notes",
			}
			runUpdateStub := mocka.Function(t, &runUpdate, int64(0), nil)
			defer runUpdateStub.Restore()

			err := UpdateDetail(tx, id, version, false, nil, values, user)

			assert.Equal(t, "transaction detail not found (42 @ 1)", err.Error())
			assert.Equal(t, sqltest.UpdateArgs(
				tx, updateTxDetailSQL, false, nil, false, nil, false, nil, true, "notes", false, nil, false, nil, user, id, version),
				runUpdateStub.GetCall(0).Arguments())
		})
	})
}

func Test_DeleteDetails(t *testing.T) {
	ids := []*VersionID{{42, 1}, {24, 0}}
	expectedErr := errors.New("delete failed")
	t.Run("returns delete transactions error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			runUpdateStub := mocka.Function(t, &runUpdate, int64(2), nil)
			runUpdateStub.OnCall(1).Return(int64(0), expectedErr)
			defer runUpdateStub.Restore()

			err := DeleteDetails(tx, ids)

			assert.Same(t, expectedErr, err)
			assert.Equal(t, 2, runUpdateStub.CallCount())
			idArg, _ := json.Marshal(ids)
			assert.Equal(t, sqltest.UpdateArgs(tx, deleteDetailsSQL, idArg), runUpdateStub.GetCall(0).Arguments())
			assert.Equal(t, sqltest.UpdateArgs(tx, deleteEmptyTransactionsSQL), runUpdateStub.GetCall(1).Arguments())
		})
	})
	t.Run("returns delete details error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			runUpdateStub := mocka.Function(t, &runUpdate, int64(0), expectedErr)
			defer runUpdateStub.Restore()

			err := DeleteDetails(tx, ids)

			assert.Same(t, expectedErr, err)
			assert.Equal(t, 1, runUpdateStub.CallCount())
		})
	})
	t.Run("returns not found error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			runUpdateStub := mocka.Function(t, &runUpdate, int64(0), nil)
			defer runUpdateStub.Restore()

			err := DeleteDetails(tx, ids)

			assert.Equal(t, "transaction detail(s) not found", err.Error())
			assert.Equal(t, 1, runUpdateStub.CallCount())
		})
	})
}

func Test_DeleteTransfer(t *testing.T) {
	relatedDetailID := int64(42)
	expectedErr := errors.New("delete failed")
	t.Run("returns delete transactions error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			runUpdateStub := mocka.Function(t, &runUpdate, int64(2), nil)
			runUpdateStub.OnCall(1).Return(int64(0), expectedErr)
			defer runUpdateStub.Restore()

			err := DeleteTransfer(tx, relatedDetailID)

			assert.Same(t, expectedErr, err)
			assert.Equal(t, 2, runUpdateStub.CallCount())
			assert.Equal(t, sqltest.UpdateArgs(tx, deleteTransferDetailSQL, relatedDetailID), runUpdateStub.GetCall(0).Arguments())
			assert.Equal(t, sqltest.UpdateArgs(tx, deleteEmptyTransactionsSQL), runUpdateStub.GetCall(1).Arguments())
		})
	})
	t.Run("returns delete detail error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			runUpdateStub := mocka.Function(t, &runUpdate, int64(0), expectedErr)
			defer runUpdateStub.Restore()

			err := DeleteTransfer(tx, relatedDetailID)

			assert.Same(t, expectedErr, err)
			assert.Equal(t, 1, runUpdateStub.CallCount())
		})
	})
	t.Run("returns no error for not found", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			runUpdateStub := mocka.Function(t, &runUpdate, int64(0), nil)
			defer runUpdateStub.Restore()

			err := DeleteTransfer(tx, relatedDetailID)

			assert.Nil(t, err)
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

			result, err := ValidateDetails(tx, transactionIDs)

			assert.Nil(t, err)
			assert.Equal(t, map[int64]string{42: "invalid shares", 96: "shares required"}, result)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})
	})
	t.Run("returns query error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			expectedErr := errors.New("query error")
			mockDB.ExpectQuery(validateDetailsSQL).WithArgs(int64sToJson(transactionIDs)).WillReturnError(expectedErr)

			_, err := ValidateDetails(tx, transactionIDs)

			assert.Same(t, expectedErr, err)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})
	})
}
