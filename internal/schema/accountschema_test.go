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

func Test_accountQueryFields_Resolve(t *testing.T) {
	companyID := int64(99)
	accountID := int64(1)
	accounts := []*domain.Account{domain.NewAccount(accountID, &companyID)}
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		getAll := mocka.Function(t, &getAllAccounts, accounts)
		byID := mocka.Function(t, &getAccountByID, accounts)
		byName := mocka.Function(t, &getAccountsByName, accounts)
		defer func() {
			getAll.Restore()
			byID.Restore()
			byName.Restore()
		}()
		name := "the account"
		tests := []struct {
			name     string
			argName  string
			argValue interface{}
			stub     *mocka.Stub
			stubArgs []interface{}
			err      bool
		}{
			{name: "returns all accounts", stub: getAll, stubArgs: []interface{}{tx}},
			{name: "returns account with ID", argName: "id", argValue: 123, stub: byID, stubArgs: []interface{}{tx, int64(123)}},
			{name: "returns accounts with name", argName: "name", argValue: name, stub: byName, stubArgs: []interface{}{tx, name}},
		}
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				params := newResolveParams(tx, accountQuery, newField("", "id"), newField("", "name")).addArg(test.argName, test.argValue)

				result, err := accountQueryFields.Resolve(params.ResolveParams)

				if test.err {
					assert.NotNil(t, err, "Expected an error")
				} else {
					assert.Nil(t, err)
					assert.Equal(t, 1, test.stub.CallCount())
					assert.Equal(t, test.stubArgs, test.stub.GetFirstCall().Arguments())
					assert.Equal(t, accounts, result)
				}
			})
		}
	})
}

type mockAccountModel struct {
	company *domain.Company
	err     error
	tx      *sql.Tx
}

func (a *mockAccountModel) GetCompany(tx *sql.Tx) (*domain.Company, error) {
	a.tx = tx
	return a.company, a.err
}

func Test_resolveCompany(t *testing.T) {
	companyID := int64(99)
	company := domain.NewCompany(companyID, "")
	mockAccount := &mockAccountModel{company: company, err: errors.New("test error")}
	sqltest.TestInTx(t, func(_ sqlmock.Sqlmock, tx *sql.Tx) {
		params := newResolveParams(tx, "", newField("", "id"), newField("", "name")).setSource(mockAccount)

		result, err := resolveCompany(params.ResolveParams)

		assert.Same(t, mockAccount.err, err)
		assert.Same(t, mockAccount.company, result)
		assert.Same(t, tx, mockAccount.tx)
	})
}
