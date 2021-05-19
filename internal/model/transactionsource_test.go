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

func Test_transactionSource_loadTransactionDetails(t *testing.T) {
	accountID := int64(42)
	txID := int64(96)
	detailID := int64(69)
	t.Run("returns existing error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			expectedErr := errors.New("test error")
			txSource := &transactionSource{err: expectedErr}

			assert.Same(t, expectedErr, txSource.loadTransactionDetails(tx))

			assert.Nil(t, mockDB.ExpectationsWereMet())
		})
	})
	t.Run("returns DAO error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			expectedErr := errors.New("test error")
			txSource := &transactionSource{accountID: accountID}
			getDetailsStub := mocka.Function(t, &getDetailsByAccountID, nil, expectedErr)
			defer getDetailsStub.Restore()

			assert.Same(t, expectedErr, txSource.loadTransactionDetails(tx))

			assert.Nil(t, mockDB.ExpectationsWereMet())
		})
	})
	t.Run("loads transaction details by account ID", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			txSource := &transactionSource{accountID: accountID}
			detail := &database.TransactionDetail{
				ID:            detailID,
				TransactionID: txID,
			}
			getDetailsStub := mocka.Function(t, &getDetailsByAccountID, []*database.TransactionDetail{detail}, nil)
			defer getDetailsStub.Restore()

			err := txSource.loadTransactionDetails(tx)

			assert.Nil(t, err)
			assert.Equal(t, detail, txSource.detailsByTxID[txID][0].TransactionDetail)
			assert.Equal(t, txSource, txSource.detailsByTxID[txID][0].txSource)
		})
	})
	t.Run("loads transaction details by tx IDs", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			txSource := &transactionSource{txIDs: []int64{txID}}
			detail := &database.TransactionDetail{
				ID:            detailID,
				TransactionID: txID,
			}
			getDetailsStub := mocka.Function(t, &getDetailsByTxIDs, []*database.TransactionDetail{detail}, nil)
			defer getDetailsStub.Restore()

			err := txSource.loadTransactionDetails(tx)

			assert.Nil(t, err)
			assert.Equal(t, detail, txSource.detailsByTxID[txID][0].TransactionDetail)
			assert.Equal(t, txSource, txSource.detailsByTxID[txID][0].txSource)
		})
	})
}

func Test_transactionSource_loadRelatedDetails(t *testing.T) {
	id := int64(42)
	txID := int64(96)
	t.Run("returns existing error", func(t *testing.T) {
		txSource := &transactionSource{err: errors.New("test error")}

		err := txSource.loadRelatedDetails(nil)

		assert.Same(t, txSource.err, err)
	})
	t.Run("returns nil for existing details", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			relatedDetail := NewTransactionDetail(id, -1)
			txSource := &transactionSource{accountID: 69, relatedDetailsByID: map[int64]*TransactionDetail{id: relatedDetail}}

			err := txSource.loadRelatedDetails(tx)

			assert.Nil(t, err)
		})
	})
	t.Run("loads related details by account ID", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			txSource := &transactionSource{accountID: 69}
			relatedDetail := &database.TransactionDetail{ID: id}
			getRelatedDetailsStub := mocka.Function(t, &getRelatedDetailsByAccountID, []*database.TransactionDetail{relatedDetail}, nil)
			defer getRelatedDetailsStub.Restore()

			err := txSource.loadRelatedDetails(tx)

			assert.Nil(t, err)
			assert.Equal(t, []interface{}{tx, txSource.accountID}, getRelatedDetailsStub.GetCall(0).Arguments())
			assert.Equal(t, relatedDetail, txSource.relatedDetailsByID[id].TransactionDetail)
			assert.Equal(t, txSource, txSource.relatedDetailsByID[id].txSource)
		})
	})
	t.Run("loads related details by transaction IDs", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			txSource := &transactionSource{txIDs: []int64{txID}}
			relatedDetail := &database.TransactionDetail{ID: id}
			getRelatedDetailsStub := mocka.Function(t, &getRelatedDetailsByTxIDs, []*database.TransactionDetail{relatedDetail}, nil)
			defer getRelatedDetailsStub.Restore()

			err := txSource.loadRelatedDetails(tx)

			assert.Nil(t, err)
			assert.Equal(t, []interface{}{tx, txSource.txIDs}, getRelatedDetailsStub.GetCall(0).Arguments())
			assert.Equal(t, relatedDetail, txSource.relatedDetailsByID[id].TransactionDetail)
			assert.Equal(t, txSource, txSource.relatedDetailsByID[id].txSource)
		})
	})
	t.Run("returns query error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			expectedErr := errors.New("test error")
			txSource := &transactionSource{accountID: 69}
			getRelatedDetailsStub := mocka.Function(t, &getRelatedDetailsByAccountID, nil, expectedErr)
			defer getRelatedDetailsStub.Restore()

			err := txSource.loadRelatedDetails(tx)

			assert.Same(t, expectedErr, err)
			assert.Same(t, expectedErr, txSource.err)
		})
	})
}

func Test_transactionSource_loadRelatedTransactions(t *testing.T) {
	id := int64(42)
	t.Run("returns existing error", func(t *testing.T) {
		txSource := &transactionSource{err: errors.New("test error")}

		err := txSource.loadRelatedTransactions(nil)

		assert.Same(t, txSource.err, err)
	})
	t.Run("loads related transactions by account id", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			txSource := &transactionSource{accountID: 69}
			transaction := &database.Transaction{ID: id}
			getTransactionsStub := mocka.Function(t, &getRelatedTransactionsByAccountID, []*database.Transaction{transaction}, nil)
			defer getTransactionsStub.Restore()

			err := txSource.loadRelatedTransactions(tx)

			assert.Nil(t, err)
			assert.Equal(t, transaction, txSource.relatedTxByID[id].Transaction)
			assert.Equal(t, txSource, txSource.relatedTxByID[id].source)
		})
	})
	t.Run("loads related transactions by transaction ids", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			txSource := &transactionSource{txIDs: []int64{69}}
			transaction := &database.Transaction{ID: id}
			getTransactionsStub := mocka.Function(t, &getRelatedTransactions, []*database.Transaction{transaction}, nil)
			defer getTransactionsStub.Restore()

			err := txSource.loadRelatedTransactions(tx)

			assert.Nil(t, err)
			assert.Equal(t, transaction, txSource.relatedTxByID[id].Transaction)
			assert.Equal(t, txSource, txSource.relatedTxByID[id].source)
		})
	})
	t.Run("returns nil for existing transactions", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			relatedTx := &Transaction{Transaction: &database.Transaction{ID: id}}
			txSource := &transactionSource{accountID: 69, relatedTxByID: map[int64]*Transaction{id: relatedTx}}

			err := txSource.loadRelatedTransactions(tx)

			assert.Nil(t, err)
		})
	})
	t.Run("returns DAO error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			expectedErr := errors.New("test error")
			txSource := &transactionSource{accountID: 69}
			getTransactionsStub := mocka.Function(t, &getRelatedTransactionsByAccountID, nil, expectedErr)
			defer getTransactionsStub.Restore()

			err := txSource.loadRelatedTransactions(tx)

			assert.Same(t, expectedErr, err)
			assert.Nil(t, txSource.relatedTxByID)
		})
	})
}
