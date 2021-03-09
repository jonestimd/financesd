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

func Test_companyQueryFields_Resolve(t *testing.T) {
	companies := []*model.Company{{ID: 1}}
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		getAll := mocka.Function(t, &getAllCompanies, companies, nil)
		byID := mocka.Function(t, &getCompanyByID, companies, nil)
		byName := mocka.Function(t, &getCompanyByName, companies, nil)
		defer func() {
			getAll.Restore()
			byID.Restore()
			byName.Restore()
		}()
		name := "the company"
		tests := []struct {
			name     string
			argName  string
			argValue interface{}
			stub     *mocka.Stub
			stubArgs []interface{}
			err      bool
		}{
			{name: "returns all companies", stub: getAll, stubArgs: []interface{}{tx}},
			{name: "returns company with ID", argName: "id", argValue: "123", stub: byID, stubArgs: []interface{}{tx, int64(123)}},
			{name: "returns company with name", argName: "name", argValue: name, stub: byName, stubArgs: []interface{}{tx, name}},
			{name: "returns error for invalid ID", argName: "id", argValue: "abc", stub: byID, err: true},
		}
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				params := newResolveParams(tx, companyQuery, newField("", "id"), newField("", "name")).addArg(test.argName, test.argValue)

				result, err := companyQueryFields().Resolve(params.ResolveParams)

				if test.err {
					assert.NotNil(t, err, "Expected an error")
				} else {
					assert.Nil(t, err, "Unexpected error: %v", err)
					assert.Equal(t, 1, test.stub.CallCount(), "Expected 1 call")
					assert.Equal(t, test.stubArgs, test.stub.GetFirstCall().Arguments(), "Args don't match")
					assert.Equal(t, companies, result, "Result doesn't match")
				}
			})
		}
	})
}

type mockCompanyModel struct {
	accounts []*model.Account
	err      error
	tx       *sql.Tx
}

func (c *mockCompanyModel) GetAccounts(tx *sql.Tx) ([]*model.Account, error) {
	c.tx = tx
	return c.accounts, c.err
}

func Test_resolveAccounts(t *testing.T) {
	companyID := int64(42)
	accounts := []*model.Account{
		{ID: 123, CompanyID: &companyID},
	}
	mockCompany := &mockCompanyModel{accounts: accounts, err: errors.New("test error")}
	sqltest.TestInTx(t, func(_ sqlmock.Sqlmock, tx *sql.Tx) {
		params := newResolveParams(tx, companyQuery, newField("", "id"), newField("", "name")).setSource(mockCompany)

		result, err := resolveAccounts(params.ResolveParams)

		assert.Same(t, mockCompany.err, err)
		assert.Equal(t, mockCompany.accounts, result)
		assert.Same(t, tx, mockCompany.tx)
	})
}

func Test_addCompanyFields_Resolve(t *testing.T) {
	names := []interface{}{"The Company"}
	companies := []*model.Company{{ID: 42, Name: names[0].(string)}}
	expectedErr := errors.New("test error")
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		mockAddCompanies := mocka.Function(t, &addCompanies, companies, expectedErr)
		defer mockAddCompanies.Restore()
		params := newResolveParams(tx, companyQuery, newField("", "id"), newField("", "name")).addArg("names", names)

		result, err := addCompaniesFields.Resolve(params.ResolveParams)

		assert.Same(t, expectedErr, err)
		assert.Equal(t, companies, result)
		assert.Equal(t, []interface{}{tx, asStrings(names), "somebody"}, mockAddCompanies.GetFirstCall().Arguments())
	})
}
