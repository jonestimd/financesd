package schema

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MonsantoCo/mocka/v2"
	"github.com/jonestimd/financesd/internal/domain"
	"github.com/jonestimd/financesd/internal/sqltest"
	"github.com/stretchr/testify/assert"
)

func Test_transactionQueryFields_Resolve_returnsRows(t *testing.T) {
	transactions := []*domain.Transaction{domain.NewTransaction(42)}
	getTransactions := mocka.Function(t, &getAccountTransactions, transactions)
	defer getTransactions.Restore()
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		accountID := 123
		params := newResolveParams(tx, transactionQuery, newField("", "id"), newField("", "memo")).addArg("accountId", accountID)

		result, err := transactionQueryFields.Resolve(params.ResolveParams)

		assert.Nil(t, err)
		assert.Equal(t, transactions, result)
		assert.Equal(t, []interface{}{tx, int64(123)}, getTransactions.GetFirstCall().Arguments())
	})
}

func Test_resolveDetails(t *testing.T) {
	txID := int64(42)
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		t.Run("returns details", func(t *testing.T) {
			details := []*domain.TransactionDetail{domain.NewTransactionDetail(96, txID)}
			domainTx := domain.NewTransaction(txID)
			domainTx.SetDetails(details)
			params := newResolveParams(tx, transactionQuery, newField("", "id")).setSource(domainTx)

			result, err := getTxSchema().Fields()["details"].Resolve(params.ResolveParams)

			assert.Nil(t, err)
			assert.Equal(t, details, result)
		})
		t.Run("returns error for invalid source", func(t *testing.T) {
			params := newResolveParams(tx, transactionQuery, newField("", "id"))

			_, err := getTxSchema().Fields()["details"].Resolve(params.ResolveParams)

			assert.Equal(t, "invalid source", err.Error())
		})
	})
}

func Test_resolveRelatedDetail(t *testing.T) {
	detailID := int64(96)
	relatedID := int64(69)
	resolver := findSchemaField(getTxSchema(), "details", "relatedDetail").Resolve
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		t.Run("returns detail", func(t *testing.T) {
			detail := domain.NewTransactionDetail(detailID, 42)
			relatedDetail := domain.NewTransactionDetail(relatedID, 24)
			detail.SetRelatedDetail(relatedDetail)
			params := newResolveParams(tx, transactionQuery, newField("", "id")).setSource(detail)

			result, err := resolver(params.ResolveParams)

			assert.Nil(t, err)
			assert.Equal(t, relatedDetail, result)
		})
		t.Run("returns error for invalid source", func(t *testing.T) {
			params := newResolveParams(tx, transactionQuery, newField("", "id"))

			_, err := resolver(params.ResolveParams)

			assert.Equal(t, "invalid source", err.Error())
		})
	})
}

func Test_resolveRelatedTransaction(t *testing.T) {
	detailID := int64(96)
	txID := int64(69)
	resolver := findSchemaField(getTxSchema(), "details", "relatedDetail", "transaction").Resolve
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		t.Run("returns transaction", func(t *testing.T) {
			detail := domain.NewTransactionDetail(detailID, txID)
			relatedTx := domain.NewTransaction(24)
			detail.SetRelatedTransaction(relatedTx)
			params := newResolveParams(tx, transactionQuery, newField("", "id")).setSource(detail)

			result, err := resolver(params.ResolveParams)

			assert.Nil(t, err)
			assert.Equal(t, relatedTx, result)
		})
		t.Run("returns error for invalid source", func(t *testing.T) {
			params := newResolveParams(tx, transactionQuery, newField("", "id"))

			_, err := resolver(params.ResolveParams)

			assert.Equal(t, "invalid source", err.Error())
		})
	})
}

func Test_updateTransactions_Resolve_delete(t *testing.T) {
	id := 42
	args := []map[string]interface{}{{"id": id, "version": 1}}
	transactions := []*domain.Transaction{}
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		deleteTransactionsStub := mocka.Function(t, &deleteTransactions)
		defer deleteTransactionsStub.Restore()
		params := newResolveParams(tx, transactionQuery, newField("", "id")).addArrayArg("delete", args)

		result, err := updateTxFields.Resolve(params.ResolveParams)

		assert.Equal(t, transactions, result)
		assert.Equal(t, []interface{}{tx, args}, deleteTransactionsStub.GetFirstCall().Arguments())
		assert.Nil(t, err)
	})
}

func Test_updateTransactions_Resolve_update(t *testing.T) {
	id := 42
	name := "new name"
	args := []map[string]interface{}{{"id": id, "name": name, "details": []map[string]interface{}{{"id": 96, "version": 1}}}}
	updateIDs := []int64{42}
	transactions := []*domain.Transaction{domain.NewTransaction(int64(id))}
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		mockUpdateTransactions := mocka.Function(t, &updateTransactions, updateIDs)
		defer mockUpdateTransactions.Restore()
		mockGetTransactions := mocka.Function(t, &getTransactionsByIDs, transactions)
		defer mockGetTransactions.Restore()
		params := newResolveParams(tx, transactionQuery, newField("", "id")).addArrayArg("update", args, "details")

		result, err := updateTxFields.Resolve(params.ResolveParams)

		assert.Equal(t, transactions, result)
		assert.Equal(t, []interface{}{tx, args, "somebody"}, mockUpdateTransactions.GetFirstCall().Arguments())
		assert.Nil(t, err)
		assert.Equal(t, []interface{}{tx, updateIDs}, mockGetTransactions.GetCall(0).Arguments())
	})
}

func Test_updateTransactions_Resolve_insert(t *testing.T) {
	name := "new name"
	accountID := 96
	args := []map[string]interface{}{{"date": "2020-12-25", "name": name, "details": []map[string]interface{}{{"amount": 96}}}}
	newIDs := []int64{42}
	transactions := []*domain.Transaction{domain.NewTransaction(newIDs[0])}
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		mockInsertTransactions := mocka.Function(t, &insertTransactions, newIDs)
		defer mockInsertTransactions.Restore()
		mockGetTransactions := mocka.Function(t, &getTransactionsByIDs, transactions)
		defer mockGetTransactions.Restore()
		params := newResolveParams(tx, transactionQuery, newField("", "id")).
			addArg("accountId", accountID).
			addArrayArg("add", args, "details")

		result, err := updateTxFields.Resolve(params.ResolveParams)

		assert.Nil(t, err)
		assert.Equal(t, transactions, result)
		assert.Equal(t, []interface{}{tx, int64(accountID), args, "somebody"}, mockInsertTransactions.GetFirstCall().Arguments())
		assert.Equal(t, []interface{}{tx, newIDs}, mockGetTransactions.GetCall(0).Arguments())
	})
}
