package model

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MonsantoCo/mocka/v2"
	"github.com/jonestimd/financesd/internal/database"
	"github.com/jonestimd/financesd/internal/sqltest"
	"github.com/stretchr/testify/assert"
)

func Test_GetTransactions(t *testing.T) {
	accountID := int64(42)
	t.Run("returns query error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			expectedErr := errors.New("test error")
			getTransactionsStub := mocka.Function(t, &getTransactions, nil, expectedErr)
			defer getTransactionsStub.Restore()

			_, err := GetTransactions(tx, accountID)

			assert.Same(t, expectedErr, err)
			assert.Equal(t, []interface{}{tx, accountID}, getTransactionsStub.GetCall(0).Arguments())
		})
	})
	t.Run("returns transactions", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			txID := int64(69)
			expectedTx := &Transaction{
				Transaction: &database.Transaction{ID: txID},
				source:      &transactionSource{accountID: accountID},
			}
			getTransactionsStub := mocka.Function(t, &getTransactions, []*database.Transaction{expectedTx.Transaction}, nil)
			defer getTransactionsStub.Restore()

			result, err := GetTransactions(tx, accountID)

			assert.Nil(t, err)
			assert.Equal(t, []*Transaction{expectedTx}, result)
		})
	})
}

func Test_GetDetails(t *testing.T) {
	txID := int64(96)
	detailID := int64(69)
	expectedErr := errors.New("test error")
	detailsByTxID := map[int64][]*TransactionDetail{txID: {NewTransactionDetail(detailID, txID)}}
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
			transaction := &Transaction{Transaction: &database.Transaction{ID: txID}, source: test.source}

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
	expectedErr := errors.New("test error")
	t.Run("update transaction", func(t *testing.T) {
		tests := []struct {
			name string
			err  error
			ids  []int64
		}{
			{"updates transactions", nil, []int64{42}},
			{"returns DAO error", expectedErr, nil},
		}
		update := database.InputObject{
			"id":      id,
			"version": version,
		}
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
					updateTransactionStub := mocka.Function(t, &updateTransaction, test.err)
					defer updateTransactionStub.Restore()

					ids, err := UpdateTransactions(tx, []map[string]interface{}{update}, user)

					assert.Equal(t, test.ids, ids)
					assert.Equal(t, test.err, err)
					assert.Equal(t, []interface{}{tx, int64(id), int64(version), update, user}, updateTransactionStub.GetCall(0).Arguments())
				})
			})
		}
	})
	t.Run("update details", func(t *testing.T) {
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
			t.Run(test.name, func(t *testing.T) {
				sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
					updateTransactionStub := mocka.Function(t, &updateTransaction, nil)
					updateTxDetailsStub := mocka.Function(t, &updateTxDetails, test.err)
					defer updateTransactionStub.Restore()
					defer updateTxDetailsStub.Restore()

					_, err := UpdateTransactions(tx, []map[string]interface{}{update}, user)

					assert.Equal(t, test.err, err)
					assert.Equal(t, []interface{}{tx, int64(id), update["details"], user}, updateTxDetailsStub.GetCall(0).Arguments())
				})
			})
		}
	})
}

func Test_GetTransactionsByIDs(t *testing.T) {
	t.Run("sets source", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			transactions := []*database.Transaction{{ID: 42}}
			getTransactionsByIDsStub := mocka.Function(t, &getTransactionsByIDs, transactions, nil)
			defer getTransactionsByIDsStub.Restore()

			result, err := GetTransactionsByIDs(tx, []int64{42})

			assert.Nil(t, err)
			assert.Equal(t, len(transactions), len(result))
			assert.Same(t, transactions[0], result[0].Transaction)
			assert.Equal(t, []int64{42}, result[0].source.txIDs)
		})
	})
	t.Run("returns query error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			expectdErr := errors.New("query error")
			getTransactionsByIDsStub := mocka.Function(t, &getTransactionsByIDs, nil, expectdErr)
			defer getTransactionsByIDsStub.Restore()

			result, err := GetTransactionsByIDs(tx, []int64{42})

			assert.Same(t, expectdErr, err)
			assert.Nil(t, result)
		})
	})
}
