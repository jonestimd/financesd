package database

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MonsantoCo/mocka/v2"
	"github.com/jonestimd/financesd/internal/database/table"
	"github.com/jonestimd/financesd/internal/sqltest"
	"github.com/stretchr/testify/assert"
)

func Test_GetTransactions(t *testing.T) {
	sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
		accountID := int64(42)
		txID := int64(69)
		expectedTx := &table.Transaction{ID: txID}
		mockDB.ExpectQuery(accountTransactionsSQL).WithArgs(accountID).WillReturnRows(sqltest.MockRows("id").AddRow(txID))

		result := GetTransactions(tx, accountID)

		assert.Equal(t, []*table.Transaction{expectedTx}, result)
		assert.Nil(t, mockDB.ExpectationsWereMet())
	})
}

func Test_GetRelatedTransactionsByAccountID(t *testing.T) {
	sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
		accountID := int64(42)
		txID := int64(69)
		expectedTx := &table.Transaction{ID: txID}
		mockDB.ExpectQuery(accountRelatedTxSQL).WithArgs(accountID).WillReturnRows(sqltest.MockRows("id").AddRow(txID))

		result := GetRelatedTransactionsByAccountID(tx, accountID)

		assert.Equal(t, []*table.Transaction{expectedTx}, result)
		assert.Nil(t, mockDB.ExpectationsWereMet())
	})
}

func Test_GetRelatedTransactions(t *testing.T) {
	sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
		relatedIDs := []int64{42}
		txID := int64(69)
		expectedTx := &table.Transaction{ID: txID}
		mockDB.ExpectQuery(relatedTxSQL).WithArgs(int64sToJson(relatedIDs)).WillReturnRows(sqltest.MockRows("id").AddRow(txID))

		result := GetRelatedTransactions(tx, relatedIDs)

		assert.Equal(t, []*table.Transaction{expectedTx}, result)
		assert.Nil(t, mockDB.ExpectationsWereMet())
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
				runUpdateStub := mocka.Function(t, &runUpdate, int64(1))
				defer runUpdateStub.Restore()

				UpdateTransaction(tx, id, version, update, user)

				assert.Equal(t, expectedArgs, runUpdateStub.GetCall(0).Arguments())
				assert.Nil(t, mockDB.ExpectationsWereMet())
			})
		})
	}
	t.Run("panics for not found", func(t *testing.T) {
		update := InputObject{}
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			runUpdateStub := mocka.Function(t, &runUpdate, int64(0))
			defer runUpdateStub.Restore()
			defer func() {
				assert.Nil(t, mockDB.ExpectationsWereMet())
				if err := recover(); err != nil {
					assert.Equal(t, "transaction not found (42 @ 1)", err.(error).Error())
				} else {
					assert.Fail(t, "expected an error")
				}
			}()

			UpdateTransaction(tx, id, version, update, user)
		})
	})
}

func Test_GetTransactionsByIDs(t *testing.T) {
	sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
		transactions := []*table.Transaction{{ID: 42}}
		runQueryStub := mocka.Function(t, &runQuery, transactions)
		defer runQueryStub.Restore()

		result := GetTransactionsByIDs(tx, []int64{42})

		assert.Equal(t, transactions, result)
	})
}

func Test_DeleteTransactions(t *testing.T) {
	ids := []map[string]interface{}{{"id": 42, "version": 1}, {"id": 24, "version": 0}}
	idArg, _ := json.Marshal(ids)
	t.Run("deletes transactions", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			runUpdateStub := mocka.Function(t, &runUpdate, int64(2))
			defer runUpdateStub.Restore()

			DeleteTransactions(tx, ids)

			assert.Equal(t,
				sqltest.UpdateArgs(tx, "delete from transaction where json_contains(?, json_object('id', id, 'version', version))", idArg),
				runUpdateStub.GetCall(0).Arguments())
		})
	})
	t.Run("panics for transactions not found", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			runUpdateStub := mocka.Function(t, &runUpdate, int64(1))
			defer runUpdateStub.Restore()
			defer func() {
				assert.Nil(t, mockDB.ExpectationsWereMet())
				if err := recover(); err != nil {
					assert.Equal(t, "transaction(s) not found", err.(error).Error())
				} else {
					assert.Fail(t, "expected an error")
				}
			}()

			DeleteTransactions(tx, ids)
		})
	})
}
