package schema

import (
	"database/sql"
	"errors"
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

type mockTxModel struct {
	details []*domain.TransactionDetail
	err     error
	tx      *sql.Tx
}

func (m *mockTxModel) GetDetails(tx *sql.Tx) ([]*domain.TransactionDetail, error) {
	m.tx = tx
	return m.details, m.err
}

func Test_resolveDetails(t *testing.T) {
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		t.Run("returns details", func(t *testing.T) {
			mockTx := &mockTxModel{details: []*domain.TransactionDetail{domain.NewTransactionDetail(42, -1)}, err: errors.New("test error")}
			params := newResolveParams(tx, transactionQuery, newField("", "id")).setSource(mockTx)

			result, err := getTxSchema().Fields()["details"].Resolve(params.ResolveParams)

			assert.Same(t, mockTx.err, err)
			assert.Equal(t, mockTx.details, result)
			assert.Same(t, tx, mockTx.tx)
		})
		t.Run("returns error for invalid source", func(t *testing.T) {
			params := newResolveParams(tx, transactionQuery, newField("", "id"))

			_, err := getTxSchema().Fields()["details"].Resolve(params.ResolveParams)

			assert.Equal(t, "invalid source", err.Error())
		})
	})
}

type mockDetailModel struct {
	relatedDetail *domain.TransactionDetail
	relatedTx     *domain.Transaction
	err           error
	tx            *sql.Tx
}

func (m *mockDetailModel) GetRelatedDetail(tx *sql.Tx) (*domain.TransactionDetail, error) {
	m.tx = tx
	return m.relatedDetail, m.err
}

func (m *mockDetailModel) GetRelatedTransaction(tx *sql.Tx) (*domain.Transaction, error) {
	m.tx = tx
	return m.relatedTx, m.err
}

func Test_resolveRelatedDetail(t *testing.T) {
	resolver := findSchemaField(getTxSchema(), "details", "relatedDetail").Resolve
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		t.Run("returns detail", func(t *testing.T) {
			mockDetail := &mockDetailModel{relatedDetail: domain.NewTransactionDetail(42, -1), err: errors.New("test error")}
			params := newResolveParams(tx, transactionQuery, newField("", "id")).setSource(mockDetail)

			result, err := resolver(params.ResolveParams)

			assert.Same(t, mockDetail.err, err)
			assert.Equal(t, mockDetail.relatedDetail, result)
			assert.Same(t, tx, mockDetail.tx)
		})
		t.Run("returns error for invalid source", func(t *testing.T) {
			params := newResolveParams(tx, transactionQuery, newField("", "id"))

			_, err := resolver(params.ResolveParams)

			assert.Equal(t, "invalid source", err.Error())
		})
	})
}

func Test_resolveRelatedTransaction(t *testing.T) {
	resolver := findSchemaField(getTxSchema(), "details", "relatedDetail", "transaction").Resolve
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		t.Run("returns transaction", func(t *testing.T) {
			mockDetail := &mockDetailModel{relatedDetail: domain.NewTransactionDetail(42, -1), err: errors.New("test error")}
			params := newResolveParams(tx, transactionQuery, newField("", "id")).setSource(mockDetail)

			result, err := resolver(params.ResolveParams)

			assert.Same(t, mockDetail.err, err)
			assert.Equal(t, mockDetail.relatedTx, result)
			assert.Same(t, tx, mockDetail.tx)
		})
		t.Run("returns error for invalid source", func(t *testing.T) {
			params := newResolveParams(tx, transactionQuery, newField("", "id"))

			_, err := resolver(params.ResolveParams)

			assert.Equal(t, "invalid source", err.Error())
		})
	})
}

func Test_updateTransactions_Resolve_update(t *testing.T) {
	id := 42
	name := "new name"
	args := []map[string]interface{}{{"id": id, "name": name}}
	updateIDs := []int64{42}
	transactions := []*domain.Transaction{domain.NewTransaction(int64(id))}
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		mockUpdateTransactions := mocka.Function(t, &updateTransactions, updateIDs)
		defer mockUpdateTransactions.Restore()
		mockGetTransactions := mocka.Function(t, &getTransactionsByIDs, transactions)
		defer mockGetTransactions.Restore()
		params := newResolveParams(tx, companyQuery, newField("", "id")).addArg("update", args)

		result, err := updateTxFields.Resolve(params.ResolveParams)

		assert.Equal(t, transactions, result)
		assert.Equal(t, []interface{}{tx, args, "somebody"}, mockUpdateTransactions.GetFirstCall().Arguments())
		assert.Nil(t, err)
		assert.Equal(t, []interface{}{tx, updateIDs}, mockGetTransactions.GetCall(0).Arguments())
	})
}
