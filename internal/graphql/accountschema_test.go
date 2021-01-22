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

func Test_accountQueryFields_Resolve(t *testing.T) {
	companyID := int64(99)
	accountID := int64(1)
	accounts := &model.Accounts{Accounts: []*model.Account{{ID: accountID, CompanyID: &companyID}}}
	getAll := mocka.Function(t, &getAllAccounts, accounts)
	byID := mocka.Function(t, &getAccountByID, accounts)
	byName := mocka.Function(t, &getAccountsByName, accounts)
	defer func() {
		getAll.Restore()
		byID.Restore()
		byName.Restore()
	}()
	name := "the account"
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		tests := []struct {
			name     string
			argName  string
			argValue interface{}
			stub     *mocka.Stub
			stubArgs []interface{}
			err      bool
		}{
			{name: "returns all accounts", stub: getAll, stubArgs: []interface{}{tx}},
			{name: "returns account with ID", argName: "id", argValue: "123", stub: byID, stubArgs: []interface{}{tx, int64(123)}},
			{name: "returns accounts with name", argName: "name", argValue: name, stub: byName, stubArgs: []interface{}{tx, name}},
			{name: "returns error for invalid ID", argName: "id", argValue: "abc", stub: byID, err: true},
		}
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				params := newResolveParams(tx, accountQuery, newField("", "id"), newField("", "name")).addArg(test.argName, test.argValue)

				result, err := accountQueryFields.Resolve(params.ResolveParams)

				if test.err {
					assert.NotNil(t, err, "Expected an error")
				} else {
					assert.Nil(t, err, "Unexpected error: %v", err)
					assert.Equal(t, 1, test.stub.CallCount())
					assert.Equal(t, test.stubArgs, test.stub.GetFirstCall().Arguments())
					assert.Equal(t, accounts.Accounts, result)
					assert.Equal(t, rootValue(params.Info)[companyIDsRootKey], []int64{companyID})
					if test.stub == getAll {
						accountsByID := map[int64]*model.Account{accountID: accounts.Accounts[0]}
						assert.Equal(t, rootValue(params.Info)[accountsRootKey], accountsByID)
					} else {
						assert.Nil(t, rootValue(params.Info)[accountsRootKey])
					}
				}
			})
		}
	})
}

func Test_resolveCompany(t *testing.T) {
	companyID := int64(99)
	company := model.Company{ID: companyID}
	companyMap := map[int64]*model.Company{companyID: &company}
	tests := []struct {
		name       string
		companyID  *int64
		result     interface{}
		companyMap interface{}
		callCount  int
	}{
		{name: "returns nil for nil company ID", callCount: 0},
		{name: "returns company from map", companyID: &companyID, result: &company, companyMap: companyMap, callCount: 0},
		{name: "returns company from database", companyID: &companyID, result: &company, callCount: 1},
	}
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				getCompanies := mocka.Function(t, &getCompaniesByIDs, companyMap, nil)
				defer getCompanies.Restore()
				account := &model.Account{CompanyID: test.companyID}
				params := newResolveParams(tx, "", newField("", "id"), newField("", "name")).
					setSource(account).
					setRootValue(companiesRootKey, test.companyMap).
					setRootValue(companyIDsRootKey, []int64{companyID})

				result, err := resolveCompany(params.ResolveParams)

				assert.Nil(t, err)
				assert.Equal(t, test.result, result)
				assert.Equal(t, test.callCount, getCompanies.CallCount())
				if account.CompanyID != nil {
					assert.Equal(t, companyMap, rootValue(params.Info)[companiesRootKey])
				}
			})
		}
	})
}

func Test_resolveCompany_returnsDatabaseError(t *testing.T) {
	companyID := int64(99)
	expectedErr := fmt.Errorf("database error")
	getCompanies := mocka.Function(t, &getCompaniesByIDs, nil, expectedErr)
	defer getCompanies.Restore()
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		account := &model.Account{CompanyID: &companyID}
		params := newResolveParams(tx, "", newField("", "id"), newField("", "name")).
			setSource(account).
			setRootValue(companyIDsRootKey, []int64{companyID})

		result, err := resolveCompany(params.ResolveParams)

		assert.Nil(t, result)
		assert.Same(t, expectedErr, err)
	})
}
