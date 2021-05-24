package domain

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MonsantoCo/mocka/v2"
	"github.com/jonestimd/financesd/internal/database"
	"github.com/jonestimd/financesd/internal/database/table"
	"github.com/jonestimd/financesd/internal/sqltest"
	"github.com/stretchr/testify/assert"
)

func Test_GetTransactions(t *testing.T) {
	accountID := int64(42)
	sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
		txID := int64(69)
		expectedTx := &Transaction{
			Transaction: &table.Transaction{ID: txID},
			source:      &transactionSource{accountID: accountID},
		}
		getTransactionsStub := mocka.Function(t, &getTransactions, []*table.Transaction{expectedTx.Transaction})
		defer getTransactionsStub.Restore()

		result := GetTransactions(tx, accountID)

		assert.Equal(t, []*Transaction{expectedTx}, result)
	})
}

func Test_GetDetails(t *testing.T) {
	txID := int64(96)
	detailID := int64(69)
	detailsByTxID := map[int64][]*TransactionDetail{txID: {NewTransactionDetail(detailID, txID)}}
	source := &transactionSource{detailsByTxID: detailsByTxID}
	transaction := &Transaction{Transaction: &table.Transaction{ID: txID}, source: source}

	result := transaction.GetDetails(nil)

	assert.Equal(t, detailsByTxID[txID], result)
}

func Test_UpdateTransactions(t *testing.T) {
	user := "user id"
	id := 42
	version := 1
	t.Run("updates transaction", func(t *testing.T) {
		update := database.InputObject{
			"id":      id,
			"version": version,
		}
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			updateTransactionStub := mocka.Function(t, &updateTransaction)
			defer updateTransactionStub.Restore()
			validateDetailsStub := mocka.Function(t, &validateDetails, nil)
			defer validateDetailsStub.Restore()

			ids := UpdateTransactions(tx, []map[string]interface{}{update}, user)

			assert.Equal(t, []int64{42}, ids)
			assert.Equal(t, []interface{}{tx, int64(id), int64(version), update, user}, updateTransactionStub.GetCall(0).Arguments())
			assert.Equal(t, []interface{}{tx, []int64{int64(id)}}, validateDetailsStub.GetCall(0).Arguments())
		})
	})
	t.Run("updates details", func(t *testing.T) {
		update := map[string]interface{}{
			"id":      id,
			"version": version,
			"details": []map[string]interface{}{},
		}
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			updateTransactionStub := mocka.Function(t, &updateTransaction)
			defer updateTransactionStub.Restore()
			updateTxDetailsStub := mocka.Function(t, &updateTxDetails)
			defer updateTxDetailsStub.Restore()
			validateDetailsStub := mocka.Function(t, &validateDetails, nil)
			defer validateDetailsStub.Restore()

			UpdateTransactions(tx, []map[string]interface{}{update}, user)

			assert.Equal(t, []interface{}{tx, int64(id), update["details"], user}, updateTxDetailsStub.GetCall(0).Arguments())
			assert.Equal(t, []interface{}{tx, []int64{int64(id)}}, validateDetailsStub.GetCall(0).Arguments())
		})
	})
}

func Test_GetTransactionsByIDs(t *testing.T) {
	sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
		transactions := []*table.Transaction{{ID: 42}}
		getTransactionsByIDsStub := mocka.Function(t, &getTransactionsByIDs, transactions)
		defer getTransactionsByIDsStub.Restore()

		result := GetTransactionsByIDs(tx, []int64{42})

		assert.Equal(t, len(transactions), len(result))
		assert.Same(t, transactions[0], result[0].Transaction)
		assert.Equal(t, []int64{42}, result[0].source.txIDs)
	})
}
