package model

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jonestimd/financesd/internal/sqltest"
	"github.com/stretchr/testify/assert"
)

func Test_transactionSource_GetDetails(t *testing.T) {
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
	t.Run("returns query error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			expectedErr := errors.New("test error")
			txSource := &transactionSource{accountID: accountID}
			mockDB.ExpectQuery(accountTxDetailsSQL).WithArgs(accountID).WillReturnError(expectedErr)

			assert.Same(t, expectedErr, txSource.loadTransactionDetails(tx))

			assert.Nil(t, mockDB.ExpectationsWereMet())
		})
	})
	t.Run("loads transaction details by account ID", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			txSource := &transactionSource{accountID: accountID}
			expectedDetail := &TransactionDetail{
				ID:            detailID,
				TransactionID: txID,
				txSource:      txSource,
			}
			mockDB.ExpectQuery(accountTxDetailsSQL).WithArgs(accountID).WillReturnRows(sqltest.MockRows("id", "transaction_id").AddRow(detailID, txID))

			err := txSource.loadTransactionDetails(tx)

			assert.Nil(t, err)
			assert.Equal(t, map[int64][]*TransactionDetail{txID: {expectedDetail}}, txSource.detailsByTxID)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})
	})
	t.Run("loads transaction details by tx IDs", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			txSource := &transactionSource{txIDs: []int64{txID}}
			expectedDetail := &TransactionDetail{
				ID:            detailID,
				TransactionID: txID,
				txSource:      txSource,
			}
			mockDB.ExpectQuery(txDetailsSQL).WithArgs(int64sToJson(txSource.txIDs)).
				WillReturnRows(sqltest.MockRows("id", "transaction_id").AddRow(detailID, txID))

			err := txSource.loadTransactionDetails(tx)

			assert.Nil(t, err)
			assert.Equal(t, map[int64][]*TransactionDetail{txID: {expectedDetail}}, txSource.detailsByTxID)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})
	})
}

func Test_transactionSource_getRelatedDetails(t *testing.T) {
	id := int64(42)
	txID := int64(96)
	t.Run("returns existing error", func(t *testing.T) {
		txSource := &transactionSource{err: errors.New("test error")}

		err := txSource.loadRelatedDetails(nil)

		assert.Same(t, txSource.err, err)
	})
	t.Run("returns nil for existing details", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			relatedDetail := &TransactionDetail{ID: id}
			txSource := &transactionSource{accountID: 69, relatedDetailsByID: map[int64]*TransactionDetail{id: relatedDetail}}

			err := txSource.loadRelatedDetails(tx)

			assert.Nil(t, err)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})
	})
	t.Run("loads related details by account ID", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			txSource := &transactionSource{accountID: 69}
			mockDB.ExpectQuery(accountRelatedDetailsSQL).WithArgs(txSource.accountID).WillReturnRows(sqltest.MockRows("id").AddRow(id))

			err := txSource.loadRelatedDetails(tx)

			assert.Nil(t, err)
			assert.Equal(t, map[int64]*TransactionDetail{id: {ID: id, txSource: txSource}}, txSource.relatedDetailsByID)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})
	})
	t.Run("loads related details by transaction IDs", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			txSource := &transactionSource{txIDs: []int64{txID}}
			mockDB.ExpectQuery(relatedDetailsSQL).WithArgs(int64sToJson(txSource.txIDs)).WillReturnRows(sqltest.MockRows("id").AddRow(id))

			err := txSource.loadRelatedDetails(tx)

			assert.Nil(t, err)
			assert.Equal(t, map[int64]*TransactionDetail{id: {ID: id, txSource: txSource}}, txSource.relatedDetailsByID)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})
	})
	t.Run("returns query error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			expectedErr := errors.New("test error")
			txSource := &transactionSource{accountID: 69}
			mockDB.ExpectQuery(accountRelatedDetailsSQL).WithArgs(txSource.accountID).WillReturnError(expectedErr)

			err := txSource.loadRelatedDetails(tx)

			assert.Same(t, expectedErr, err)
			assert.Same(t, expectedErr, txSource.err)
			assert.Nil(t, mockDB.ExpectationsWereMet())
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
	t.Run("loads related transactions", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			txSource := &transactionSource{accountID: 69}
			mockDB.ExpectQuery(accountRelatedTxSQL).WithArgs(txSource.accountID).WillReturnRows(sqltest.MockRows("id").AddRow(id))

			err := txSource.loadRelatedTransactions(tx)

			assert.Nil(t, err)
			assert.Equal(t, map[int64]*Transaction{id: {ID: id}}, txSource.relatedTxByID)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})
	})
	t.Run("returns nil for existing transactions", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			relatedTx := &Transaction{ID: id}
			txSource := &transactionSource{accountID: 69, relatedTxByID: map[int64]*Transaction{id: relatedTx}}

			err := txSource.loadRelatedTransactions(tx)

			assert.Nil(t, err)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})
	})
	t.Run("returns query error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			expectedErr := errors.New("test error")
			txSource := &transactionSource{accountID: 69}
			mockDB.ExpectQuery(accountRelatedTxSQL).WithArgs(txSource.accountID).WillReturnError(expectedErr)

			err := txSource.loadRelatedTransactions(tx)

			assert.Same(t, expectedErr, err)
			assert.Nil(t, txSource.relatedTxByID)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})
	})
}
