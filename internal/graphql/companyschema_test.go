package graphql

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MonsantoCo/mocka/v2"
	"github.com/jonestimd/financesd/internal/model"
	"github.com/jonestimd/financesd/internal/sqltest"
	"github.com/stretchr/testify/assert"
)

func Test_companyQueryFields_Resolve(t *testing.T) {
	companies := &model.Companies{Companies: []*model.Company{{ID: 1}}}
	getAll := mocka.Function(t, &getAllCompanies, companies)
	byID := mocka.Function(t, &getCompanyByID, companies)
	byName := mocka.Function(t, &getCompanyByName, companies)
	defer func() {
		getAll.Restore()
		byID.Restore()
		byName.Restore()
	}()
	name := "the company"
	sqltest.TestQuery(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
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

				result, err := companyQueryFields.Resolve(params.ResolveParams)

				if test.err {
					assert.NotNil(t, err, "Expected an error")
				} else {
					assert.Nil(t, err, "Unexpected error: %v", err)
					assert.Equal(t, 1, test.stub.CallCount(), "Expected 1 call")
					assert.Equal(t, test.stubArgs, test.stub.GetFirstCall().Arguments(), "Args don't match")
					assert.Equal(t, companies.Companies, result, "Result doesn't match")
					assert.Equal(t, rootValue(params.Info)[companyIDsRootKey], []int64{companies.Companies[0].ID})
				}
			})
		}
	})
}

func Test_resolveAccounts(t *testing.T) {
	company := model.Company{ID: 99}
	otherCompanyID := int64(66)
	accounts := model.Accounts{Accounts: []*model.Account{
		{ID: 123, CompanyID: &company.ID},
		{ID: 456},
		{ID: 789, CompanyID: &otherCompanyID},
	}}
	expectedAccounts := []*model.Account{accounts.Accounts[0]}
	tests := []struct {
		name       string
		result     interface{}
		accountMap interface{}
		callCount  int
	}{
		{name: "returns empty array for no accounts", result: []*model.Account{}, accountMap: map[int64]*model.Account{}, callCount: 0},
		{name: "returns accounts from map", result: expectedAccounts, accountMap: accounts.ByID(), callCount: 0},
		{name: "returns accounts from database", result: expectedAccounts, callCount: 1},
	}
	sqltest.TestQuery(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				getAccounts := mocka.Function(t, &getAccountsByCompanyIDs, &accounts)
				defer getAccounts.Restore()
				params := newResolveParams(tx, companyQuery, newField("", "id"), newField("", "name")).
					setSource(&company).
					setRootValue(companyIDsRootKey, []int64{company.ID}).
					setRootValue(accountsRootKey, test.accountMap)

				result, err := resolveAccounts(params.ResolveParams)

				assert.Nil(t, err)
				assert.Equal(t, test.result, result)
				if test.accountMap != nil && len(test.accountMap.(map[int64]*model.Account)) > 0 {
					assert.Equal(t, accounts.ByID(), rootValue(params.Info)[accountsRootKey])
				}
			})
		}
	})
}

func Test_resolveAccounts_returnsDatabaseError(t *testing.T) {
	company := model.Company{ID: 99}
	expectedErr := fmt.Errorf("database error")
	getAccounts := mocka.Function(t, &getAccountsByCompanyIDs, model.NewAccounts(expectedErr))
	defer getAccounts.Restore()
	sqltest.TestQuery(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		params := newResolveParams(tx, companyQuery, newField("", "id"), newField("", "name")).
			setSource(&company).
			setRootValue(companyIDsRootKey, []int64{company.ID})

		result, err := resolveAccounts(params.ResolveParams)

		assert.Same(t, expectedErr, err)
		assert.Nil(t, result)
	})
}

func Test_resolveAccounts_returnsErrorForNoCompanyIDs(t *testing.T) {
	company := model.Company{ID: 99}
	getAccounts := mocka.Function(t, &getAccountsByCompanyIDs, nil)
	defer getAccounts.Restore()
	sqltest.TestQuery(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		params := newResolveParams(tx, companyQuery, newField("", "id"), newField("", "name")).setSource(&company)

		result, err := resolveAccounts(params.ResolveParams)

		assert.NotNil(t, err)
		assert.Nil(t, result)
		assert.Equal(t, 0, getAccounts.CallCount())
	})
}
