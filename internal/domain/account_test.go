package domain

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MonsantoCo/mocka/v2"
	"github.com/jonestimd/financesd/internal/database/table"
	"github.com/jonestimd/financesd/internal/sqltest"
	"github.com/stretchr/testify/assert"
)

func Test_Account_GetCompany(t *testing.T) {
	companyID := int64(42)
	company := &Company{Company: &table.Company{ID: companyID}}
	account := &table.Account{CompanyID: &companyID}
	loadedSource := &companySource{companiesByID: map[int64]*Company{companyID: company}}
	tests := []struct {
		name    string
		account *Account
		company *Company
		err     error
	}{
		{"returns nil for nil company ID", &Account{Account: &table.Account{}}, nil, nil},
		{"returns existing source company", &Account{Account: account, source: loadedSource}, company, nil},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.account.GetCompany(nil)

			assert.Equal(t, test.company, result)
		})
	}
}

func Test_GetAllAccounts(t *testing.T) {
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		dbAccounts := []*table.Account{{ID: 1}}
		getAllAccountsStub := mocka.Function(t, &getAllAccounts, dbAccounts)
		defer getAllAccountsStub.Restore()

		result := GetAllAccounts(tx)

		assert.Equal(t, []interface{}{tx}, getAllAccountsStub.GetFirstCall().Arguments())
		assert.Equal(t, dbAccounts[0], result[0].Account)
		assert.NotNil(t, result[0].source)
	})
}

func Test_GetAccountByID(t *testing.T) {
	id := int64(42)
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		dbAccounts := []*table.Account{{ID: 1}}
		getAccountByIDStub := mocka.Function(t, &getAccountByID, dbAccounts)
		defer getAccountByIDStub.Restore()

		result := GetAccountByID(tx, id)

		assert.Equal(t, []interface{}{tx, id}, getAccountByIDStub.GetFirstCall().Arguments())
		assert.Equal(t, dbAccounts[0], result[0].Account)
		assert.NotNil(t, result[0].source)
	})
}

func Test_GetAccountsByName(t *testing.T) {
	name := "account name"
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		dbAccounts := []*table.Account{{ID: 1}}
		getAccountsByNameStub := mocka.Function(t, &getAccountsByName, dbAccounts)
		defer getAccountsByNameStub.Restore()

		result := GetAccountsByName(tx, name)

		assert.Equal(t, []interface{}{tx, name}, getAccountsByNameStub.GetFirstCall().Arguments())
		assert.Equal(t, dbAccounts[0], result[0].Account)
		assert.NotNil(t, result[0].source)
	})
}

func Test_GetAccountsByCompanyID(t *testing.T) {
	ids := []int64{42, 69}
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		dbAccounts := []*table.Account{{ID: 1}}
		getAccountsByCompanyIDsStub := mocka.Function(t, &getAccountsByCompanyIDs, dbAccounts)
		defer getAccountsByCompanyIDsStub.Restore()

		result := GetAccountsByCompanyIDs(tx, ids)

		assert.Equal(t, []interface{}{tx, ids}, getAccountsByCompanyIDsStub.GetFirstCall().Arguments())
		assert.Equal(t, dbAccounts[0], result[0].Account)
		assert.NotNil(t, result[0].source)
	})
}
