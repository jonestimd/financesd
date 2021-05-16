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

func Test_Transaction_ptrTo(t *testing.T) {
	transaction := &Transaction{}
	tests := []struct {
		column string
		ptr    interface{}
	}{
		{column: "id", ptr: &transaction.ID},
		{column: "date", ptr: &transaction.Date},
		{column: "memo", ptr: &transaction.Memo},
		{column: "reference_number", ptr: &transaction.ReferenceNumber},
		{column: "cleared", ptr: &transaction.Cleared},
		{column: "account_id", ptr: &transaction.AccountID},
		{column: "payee_id", ptr: &transaction.PayeeID},
		{column: "security_id", ptr: &transaction.SecurityID},
		{column: "version", ptr: &transaction.Version},
		{column: "change_user", ptr: &transaction.ChangeUser},
		{column: "change_date", ptr: &transaction.ChangeDate},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s returns pointer to field", test.column), func(t *testing.T) {
			field := transaction.ptrTo(test.column)

			assert.Same(t, test.ptr, field)
		})
	}
}

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
			expectedTx := &Transaction{
				ID:     txID,
				source: &transactionSource{accountID: accountID},
			}
			mockDB.ExpectQuery(accountTransactionsSQL).WithArgs(accountID).WillReturnRows(sqltest.MockRows("id").AddRow(txID))

			result, err := GetTransactions(tx, accountID)

			assert.Nil(t, err)
			assert.Equal(t, []*Transaction{expectedTx}, result)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})
	})
}

func Test_GetDetails(t *testing.T) {
	txID := int64(96)
	detailID := int64(69)
	expectedErr := errors.New("test error")
	detailsByTxID := map[int64][]*TransactionDetail{txID: {&TransactionDetail{ID: detailID, TransactionID: txID}}}
	tests := []struct {
		name           string
		source         *transactionSource
		expectedErr    error
		expectedResult []*TransactionDetail
	}{
		{"returns existing error", &transactionSource{err: expectedErr}, expectedErr, nil},
		{"returns existing details", &transactionSource{detailsByTxID: detailsByTxID}, nil, detailsByTxID[txID]},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			transaction := &Transaction{ID: txID, source: test.source}

			result, err := transaction.GetDetails(nil)

			assert.Equal(t, test.expectedErr, err)
			assert.Equal(t, test.expectedResult, result)
		})
	}
}

func Test_UpdateTransactions(t *testing.T) {
	user := "user id"
	id := 42
	version := 1
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
			update := inputObject{
				"id":      id,
				"version": version,
			}
			update[test.field] = test.value
			sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
				expectedArgs := []interface{}{
					tx, updateTxSQL,
					[]interface{}{update["date"],
						test.field == "referenceNumber", update["referenceNumber"],
						test.field == "payeeId", update.intOrNull("payeeId"),
						test.field == "securityId", update.intOrNull("securityId"),
						test.field == "memo", update["memo"],
						update["cleared"],
						test.field == "accountId", update.intOrNull("accountId"),
						user, int64(id), int64(version)}}
				runUpdateStub := mocka.Function(t, &runUpdate, int64(1), nil)
				defer runUpdateStub.Restore()

				ids, err := UpdateTransactions(tx, []map[string]interface{}{update}, user)

				assert.Nil(t, err)
				assert.Equal(t, []int64{42}, ids)
				assert.Equal(t, expectedArgs, runUpdateStub.GetCall(0).Arguments())
				assert.Nil(t, mockDB.ExpectationsWereMet())
			})
		})
	}
	expectedErr := errors.New("test error")
	t.Run("updates details", func(t *testing.T) {
		update := map[string]interface{}{
			"id":      id,
			"version": version,
			"details": []map[string]interface{}{},
		}
		tests := []struct {
			name string
			err  error
		}{
			{"updates details", nil},
			{"returns update detail error", expectedErr},
		}
		for _, test := range tests {
			sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
				runUpdateStub := mocka.Function(t, &runUpdate, int64(1), nil)
				updateTxDetailsStub := mocka.Function(t, &updateTxDetails, test.err)
				defer runUpdateStub.Restore()
				defer updateTxDetailsStub.Restore()

				_, err := UpdateTransactions(tx, []map[string]interface{}{update}, user)

				assert.Equal(t, test.err, err)
				assert.Equal(t, []interface{}{tx, int64(id), update["details"], user}, updateTxDetailsStub.GetCall(0).Arguments())
			})
		}
	})
	t.Run("returns query error", func(t *testing.T) {
		update := map[string]interface{}{
			"id":      id,
			"version": version,
		}
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			runUpdateStub := mocka.Function(t, &runUpdate, int64(0), expectedErr)
			defer runUpdateStub.Restore()

			ids, err := UpdateTransactions(tx, []map[string]interface{}{update}, user)

			assert.Same(t, expectedErr, err)
			assert.Nil(t, ids)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})
	})
	t.Run("returns error for not found", func(t *testing.T) {
		update := map[string]interface{}{
			"id":      id,
			"version": version,
		}
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			runUpdateStub := mocka.Function(t, &runUpdate, int64(0), nil)
			defer runUpdateStub.Restore()

			ids, err := UpdateTransactions(tx, []map[string]interface{}{update}, user)

			assert.Equal(t, "transaction not found (42 @ 1)", err.Error())
			assert.Nil(t, ids)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})
	})
}

func Test_GetTransactionsByIDs(t *testing.T) {
	t.Run("sets source", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			transactions := []*Transaction{{ID: 42}}
			runQueryStub := mocka.Function(t, &runQuery, transactions, nil)
			defer runQueryStub.Restore()

			result, err := GetTransactionsByIDs(tx, []int64{42})

			assert.Nil(t, err)
			assert.Equal(t, transactions, result)
			assert.Equal(t, []int64{42}, transactions[0].source.txIDs)
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
