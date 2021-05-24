package domain

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MonsantoCo/mocka/v2"
	"github.com/jonestimd/financesd/internal/database/table"
	"github.com/jonestimd/financesd/internal/sqltest"
	"github.com/stretchr/testify/assert"
)

func Test_transactionSource_loadTransactionDetails(t *testing.T) {
	accountID := int64(42)
	txID := int64(96)
	detailID := int64(69)
	t.Run("loads transaction details by account ID", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			txSource := &transactionSource{accountID: accountID}
			detail := &table.TransactionDetail{
				ID:            detailID,
				TransactionID: txID,
			}
			getDetailsStub := mocka.Function(t, &getDetailsByAccountID, []*table.TransactionDetail{detail})
			defer getDetailsStub.Restore()

			txSource.loadTransactionDetails(tx)

			assert.Equal(t, detail, txSource.detailsByTxID[txID][0].TransactionDetail)
			assert.Equal(t, txSource, txSource.detailsByTxID[txID][0].txSource)
		})
	})
	t.Run("loads transaction details by tx IDs", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			txSource := &transactionSource{txIDs: []int64{txID}}
			detail := &table.TransactionDetail{
				ID:            detailID,
				TransactionID: txID,
			}
			getDetailsStub := mocka.Function(t, &getDetailsByTxIDs, []*table.TransactionDetail{detail})
			defer getDetailsStub.Restore()

			txSource.loadTransactionDetails(tx)

			assert.Equal(t, detail, txSource.detailsByTxID[txID][0].TransactionDetail)
			assert.Equal(t, txSource, txSource.detailsByTxID[txID][0].txSource)
		})
	})
}

func Test_transactionSource_loadRelatedDetails(t *testing.T) {
	id := int64(42)
	txID := int64(96)
	t.Run("returns for existing details", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			relatedDetail := NewTransactionDetail(id, -1)
			txSource := &transactionSource{accountID: 69, relatedDetailsByID: map[int64]*TransactionDetail{id: relatedDetail}}

			txSource.loadRelatedDetails(tx)
		})
	})
	t.Run("loads related details by account ID", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			txSource := &transactionSource{accountID: 69}
			relatedDetail := &table.TransactionDetail{ID: id}
			getRelatedDetailsStub := mocka.Function(t, &getRelatedDetailsByAccountID, []*table.TransactionDetail{relatedDetail})
			defer getRelatedDetailsStub.Restore()

			txSource.loadRelatedDetails(tx)

			assert.Equal(t, []interface{}{tx, txSource.accountID}, getRelatedDetailsStub.GetCall(0).Arguments())
			assert.Equal(t, relatedDetail, txSource.relatedDetailsByID[id].TransactionDetail)
			assert.Equal(t, txSource, txSource.relatedDetailsByID[id].txSource)
		})
	})
	t.Run("loads related details by transaction IDs", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			txSource := &transactionSource{txIDs: []int64{txID}}
			relatedDetail := &table.TransactionDetail{ID: id}
			getRelatedDetailsStub := mocka.Function(t, &getRelatedDetailsByTxIDs, []*table.TransactionDetail{relatedDetail})
			defer getRelatedDetailsStub.Restore()

			txSource.loadRelatedDetails(tx)

			assert.Equal(t, []interface{}{tx, txSource.txIDs}, getRelatedDetailsStub.GetCall(0).Arguments())
			assert.Equal(t, relatedDetail, txSource.relatedDetailsByID[id].TransactionDetail)
			assert.Equal(t, txSource, txSource.relatedDetailsByID[id].txSource)
		})
	})
}

func Test_transactionSource_loadRelatedTransactions(t *testing.T) {
	id := int64(42)
	t.Run("loads related transactions by account id", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			txSource := &transactionSource{accountID: 69}
			transaction := &table.Transaction{ID: id}
			getTransactionsStub := mocka.Function(t, &getRelatedTransactionsByAccountID, []*table.Transaction{transaction})
			defer getTransactionsStub.Restore()

			txSource.loadRelatedTransactions(tx)

			assert.Equal(t, transaction, txSource.relatedTxByID[id].Transaction)
			assert.Equal(t, txSource, txSource.relatedTxByID[id].source)
		})
	})
	t.Run("loads related transactions by transaction ids", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			txSource := &transactionSource{txIDs: []int64{69}}
			transaction := &table.Transaction{ID: id}
			getTransactionsStub := mocka.Function(t, &getRelatedTransactions, []*table.Transaction{transaction})
			defer getTransactionsStub.Restore()

			txSource.loadRelatedTransactions(tx)

			assert.Equal(t, transaction, txSource.relatedTxByID[id].Transaction)
			assert.Equal(t, txSource, txSource.relatedTxByID[id].source)
		})
	})
	t.Run("returns for existing transactions", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			relatedTx := &Transaction{Transaction: &table.Transaction{ID: id}}
			txSource := &transactionSource{accountID: 69, relatedTxByID: map[int64]*Transaction{id: relatedTx}}

			txSource.loadRelatedTransactions(tx)
		})
	})
}
