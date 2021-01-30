package schema

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MonsantoCo/mocka/v2"
	"github.com/jonestimd/financesd/internal/model"
	"github.com/jonestimd/financesd/internal/sqltest"
	"github.com/stretchr/testify/assert"
)

func Test_transactionQueryFields_Resolve_returnsRows(t *testing.T) {
	transactions := []*model.Transaction{{ID: 42}}
	getTransactions := mocka.Function(t, &getAccountTransactions, transactions, nil)
	defer getTransactions.Restore()
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		accountID := "123"
		params := newResolveParams(tx, transactionQuery, newField("", "id"), newField("", "memo")).addArg("accountId", accountID)

		result, err := transactionQueryFields.Resolve(params.ResolveParams)

		assert.Nil(t, err)
		assert.Equal(t, transactions, result)
		assert.Equal(t, []interface{}{tx, int64(123)}, getTransactions.GetFirstCall().Arguments())
	})
}

func Test_transactionQueryFields_Resolve_requiresAccountID(t *testing.T) {
	getTransactions := mocka.Function(t, &getAccountTransactions, nil, nil)
	defer getTransactions.Restore()
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		params := newResolveParams(tx, transactionQuery, newField("", "id"), newField("", "memo")).addArg("accountId", "abc")

		_, err := transactionQueryFields.Resolve(params.ResolveParams)

		assert.NotNil(t, err)
		assert.Equal(t, 0, getTransactions.CallCount())
	})
}

type mockTxModel struct {
	details []*model.TransactionDetail
	err     error
	tx      *sql.Tx
}

func (m *mockTxModel) GetDetails(tx *sql.Tx) ([]*model.TransactionDetail, error) {
	m.tx = tx
	return m.details, m.err
}

func Test_resolveDetails(t *testing.T) {
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		t.Run("returns details", func(t *testing.T) {
			mockTx := &mockTxModel{details: []*model.TransactionDetail{{ID: 42}}, err: errors.New("test error")}
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
	relatedDetail *model.TransactionDetail
	relatedTx     *model.Transaction
	err           error
	tx            *sql.Tx
}

func (m *mockDetailModel) GetRelatedDetail(tx *sql.Tx) (*model.TransactionDetail, error) {
	m.tx = tx
	return m.relatedDetail, m.err
}

func (m *mockDetailModel) GetRelatedTransaction(tx *sql.Tx) (*model.Transaction, error) {
	m.tx = tx
	return m.relatedTx, m.err
}

func Test_resolveRelatedDetail(t *testing.T) {
	resolver := findSchemaField(getTxSchema(), "details", "relatedDetail").Resolve
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		t.Run("returns detail", func(t *testing.T) {
			mockDetail := &mockDetailModel{relatedDetail: &model.TransactionDetail{ID: 42}, err: errors.New("test error")}
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
			mockDetail := &mockDetailModel{relatedDetail: &model.TransactionDetail{ID: 42}, err: errors.New("test error")}
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
