package database

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MonsantoCo/mocka/v2"
	"github.com/jonestimd/financesd/internal/sqltest"
	"github.com/stretchr/testify/assert"
)

func Test_GetTransactions(t *testing.T) {
	accountID := int64(42)
	t.Run("returns query error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			expectedErr := errors.New("test error")
			mockDB.ExpectQuery(accountTransactionsSQL).WithArgs(accountID).WillReturnError(expectedErr)

			_, err := GetTransactions(tx, accountID)

			assert.Same(t, expectedErr, err)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})
	})
	t.Run("returns transactions", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			txID := int64(69)
			expectedTx := &Transaction{ID: txID}
			mockDB.ExpectQuery(accountTransactionsSQL).WithArgs(accountID).WillReturnRows(sqltest.MockRows("id").AddRow(txID))

			result, err := GetTransactions(tx, accountID)

			assert.Nil(t, err)
			assert.Equal(t, []*Transaction{expectedTx}, result)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})
	})
}

func Test_UpdateTransaction(t *testing.T) {
	user := "user id"
	id := int64(42)
	version := int64(1)
	tests := []struct {
		name  string
		field string
		value interface{}
	}{
		{"sets date", "date", "2020-12-25"},
		{"sets reference number", "referenceNumber", "ref #"},
		{"sets payee", "payeeId", 96},
		{"sets security", "securityId", 96},
		{"sets memo", "memo", "notes"},
		{"sets cleared", "cleared", "Y"},
		{"sets account", "accountId", 96},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			update := InputObject{}
			update[test.field] = test.value
			sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
				expectedArgs := []interface{}{
					tx, updateTxSQL,
					[]interface{}{update["date"],
						test.field == "referenceNumber", update["referenceNumber"],
						test.field == "payeeId", update.IntOrNull("payeeId"),
						test.field == "securityId", update.IntOrNull("securityId"),
						test.field == "memo", update["memo"],
						update["cleared"],
						test.field == "accountId", update.IntOrNull("accountId"),
						user, id, version}}
				runUpdateStub := mocka.Function(t, &runUpdate, int64(1), nil)
				defer runUpdateStub.Restore()

				err := UpdateTransaction(tx, id, version, update, user)

				assert.Nil(t, err)
				assert.Equal(t, expectedArgs, runUpdateStub.GetCall(0).Arguments())
				assert.Nil(t, mockDB.ExpectationsWereMet())
			})
		})
	}
	expectedErr := errors.New("test error")
	t.Run("returns query error", func(t *testing.T) {
		update := InputObject{}
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			runUpdateStub := mocka.Function(t, &runUpdate, int64(0), expectedErr)
			defer runUpdateStub.Restore()

			err := UpdateTransaction(tx, id, version, update, user)

			assert.Same(t, expectedErr, err)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})
	})
	t.Run("returns error for not found", func(t *testing.T) {
		update := InputObject{}
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			runUpdateStub := mocka.Function(t, &runUpdate, int64(0), nil)
			defer runUpdateStub.Restore()

			err := UpdateTransaction(tx, id, version, update, user)

			assert.Equal(t, "transaction not found (42 @ 1)", err.Error())
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})
	})
}

func Test_GetTransactionsByIDs(t *testing.T) {
	t.Run("returns transactions", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			transactions := []*Transaction{{ID: 42}}
			runQueryStub := mocka.Function(t, &runQuery, transactions, nil)
			defer runQueryStub.Restore()

			result, err := GetTransactionsByIDs(tx, []int64{42})

			assert.Nil(t, err)
			assert.Equal(t, transactions, result)
		})
	})
	t.Run("returns query error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			expectdErr := errors.New("query error")
			runQueryStub := mocka.Function(t, &runQuery, nil, expectdErr)
			defer runQueryStub.Restore()

			result, err := GetTransactionsByIDs(tx, []int64{42})

			assert.Same(t, expectdErr, err)
			assert.Nil(t, result)
		})
	})
}
