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

func Test_companyQueryFields_Resolve(t *testing.T) {
	companies := []*domain.Company{domain.NewCompany(1, "")}
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		getAll := mocka.Function(t, &getAllCompanies, companies)
		byID := mocka.Function(t, &getCompanyByID, companies)
		byName := mocka.Function(t, &getCompanyByName, companies)
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
			{name: "returns company with ID", argName: "id", argValue: 123, stub: byID, stubArgs: []interface{}{tx, int64(123)}},
			{name: "returns company with name", argName: "name", argValue: name, stub: byName, stubArgs: []interface{}{tx, name}},
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
	accounts []*domain.Account
	err      error
	tx       *sql.Tx
}

func (c *mockCompanyModel) GetAccounts(tx *sql.Tx) ([]*domain.Account, error) {
	c.tx = tx
	return c.accounts, c.err
}

func Test_resolveAccounts(t *testing.T) {
	companyID := int64(42)
	accounts := []*domain.Account{domain.NewAccount(123, &companyID)}
	mockCompany := &mockCompanyModel{accounts: accounts, err: errors.New("test error")}
	sqltest.TestInTx(t, func(_ sqlmock.Sqlmock, tx *sql.Tx) {
		params := newResolveParams(tx, companyQuery, newField("", "id"), newField("", "name")).setSource(mockCompany)

		result, err := resolveAccounts(params.ResolveParams)

		assert.Same(t, mockCompany.err, err)
		assert.Equal(t, mockCompany.accounts, result)
		assert.Same(t, tx, mockCompany.tx)
	})
}

func Test_updateCompany_Resolve_add(t *testing.T) {
	names := []interface{}{"The Company"}
	companies := []*domain.Company{domain.NewCompany(42, names[0].(string))}
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		mockAddCompanies := mocka.Function(t, &addCompanies, companies)
		defer mockAddCompanies.Restore()
		params := newResolveParams(tx, companyQuery, newField("", "id"), newField("", "name")).addArg("add", names)

		result, err := updateCompaniesFields.Resolve(params.ResolveParams)

		assert.Nil(t, err)
		assert.Equal(t, companies, result)
		assert.Equal(t, []interface{}{tx, asStrings(names), "somebody"}, mockAddCompanies.GetFirstCall().Arguments())
	})
}

func Test_updateCompany_Resolve_delete(t *testing.T) {
	ids := []interface{}{1, 3}
	companies := []*domain.Company{}
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		count := int64(2)
		mockDeleteCompanies := mocka.Function(t, &deleteCompanies, count)
		defer mockDeleteCompanies.Restore()
		params := newResolveParams(tx, companyQuery, newField("", "id")).addArg("delete", ids)

		result, err := updateCompaniesFields.Resolve(params.ResolveParams)

		assert.Nil(t, err)
		assert.Equal(t, companies, result)
		assert.Equal(t, []interface{}{tx, asInts(ids)}, mockDeleteCompanies.GetFirstCall().Arguments())
	})
}

func Test_updateCompany_Resolve_update(t *testing.T) {
	id := 42
	name := "new name"
	args := []interface{}{map[string]interface{}{"id": id, "name": name}}
	companies := []*domain.Company{}
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		mockUpdateCompanies := mocka.Function(t, &updateCompanies, companies)
		defer mockUpdateCompanies.Restore()
		params := newResolveParams(tx, companyQuery, newField("", "id")).addArg("update", args)

		result, err := updateCompaniesFields.Resolve(params.ResolveParams)

		assert.Nil(t, err)
		assert.Equal(t, companies, result)
		assert.Equal(t, []interface{}{tx, args, "somebody"}, mockUpdateCompanies.GetFirstCall().Arguments())
	})
}
