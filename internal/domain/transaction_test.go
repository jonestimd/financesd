package domain

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MonsantoCo/mocka/v2"
	"github.com/graphql-go/graphql"
	"github.com/jonestimd/financesd/internal/database"
	"github.com/jonestimd/financesd/internal/database/table"
	"github.com/jonestimd/financesd/internal/sqltest"
	"github.com/stretchr/testify/assert"
)

func Test_Transaction_Resolve(t *testing.T) {
	expectedResult := "result"
	defaultResolveStub := mocka.Function(t, &defaultResolveFn, expectedResult, nil)
	defer defaultResolveStub.Restore()
	p := graphql.ResolveParams{}
	transaction := NewTransaction(42)

	result, err := transaction.Resolve(p)

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
	assert.Same(t, transaction.Transaction, defaultResolveStub.GetCall(0).Arguments()[0].(graphql.ResolveParams).Source)
}

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

func Test_InsertTransactions(t *testing.T) {
	accountID := int64(42)
	user := "user id"
	txID := int64(96)
	errTests := []struct {
		name    string
		inserts []map[string]interface{}
		err     string
	}{
		{"panics for no details", []map[string]interface{}{{"date": "2020-12-25"}}, "new transaction requires at least 1 detail"},
		{"panics for empty details", []map[string]interface{}{{"details": []map[string]interface{}{}}}, "new transaction requires at least 1 detail"},
		{"panics for no amount", []map[string]interface{}{{"details": []map[string]interface{}{{"memo": "x"}}}}, "new transaction detail requires amount"},
	}
	for _, test := range errTests {
		t.Run(test.name, func(t *testing.T) {
			sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
				insertTransactionStub := mocka.Function(t, &insertTransaction, txID)
				defer func() {
					insertTransactionStub.Restore()
					if err := recover(); err != nil {
						assert.Equal(t, test.err, err.(error).Error())
					} else {
						assert.Fail(t, "expected an error")
					}
				}()

				InsertTransactions(tx, accountID, test.inserts, user)
			})
		})
	}
	t.Run("returns ids", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			amount := 234.56
			details := []map[string]interface{}{{"amount": amount}}
			inserts := []map[string]interface{}{
				{"date": "2020-12-25", "details": details},
			}
			insertTransactionStub := mocka.Function(t, &insertTransaction, txID)
			defer insertTransactionStub.Restore()
			insertDetailStub := mocka.Function(t, &insertDetail)
			defer insertDetailStub.Restore()
			validateDetailsStub := mocka.Function(t, &validateDetails)
			defer validateDetailsStub.Restore()

			result := InsertTransactions(tx, accountID, inserts, user)

			assert.Equal(t, []int64{txID}, result)
			assert.Equal(t, []interface{}{tx, accountID, database.InputObject(inserts[0]), user}, insertTransactionStub.GetCall(0).Arguments())
			assert.Equal(t, []interface{}{tx, txID, amount, database.InputObject(details[0]), user}, insertDetailStub.GetCall(0).Arguments())
			assert.Equal(t, []interface{}{tx, []int64{txID}}, validateDetailsStub.GetCall(0).Arguments())
		})
	})
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
			validateDetailsStub := mocka.Function(t, &validateDetails)
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
			validateDetailsStub := mocka.Function(t, &validateDetails)
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

func Test_DeleteTransactions(t *testing.T) {
	txIDs := []map[string]interface{}{{"id": 1, "version": 0}, {"id": 2, "version": 9}}
	sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
		deleteRelatedDetailsStub := mocka.Function(t, &deleteRelatedDetails)
		defer deleteRelatedDetailsStub.Restore()
		deleteTransactionDetailsStub := mocka.Function(t, &deleteTransactionDetails)
		defer deleteTransactionDetailsStub.Restore()
		deleteTransactionsStub := mocka.Function(t, &deleteTransactions)
		defer deleteTransactionsStub.Restore()

		DeleteTransactions(tx, txIDs)

		assert.Equal(t, 1, deleteRelatedDetailsStub.CallCount())
		assert.Equal(t, 1, deleteTransactionDetailsStub.CallCount())
		assert.Equal(t, 1, deleteTransactionsStub.CallCount())
	})
}
