package domain

import (
	"database/sql"
	"errors"
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
	expectedErr := errors.New("test error")
	tests := []struct {
		name    string
		account *Account
		company *Company
		err     error
	}{
		{"returns nil for nil company ID", &Account{Account: &table.Account{}}, nil, nil},
		{"returns existing source error", &Account{Account: account, source: &companySource{err: expectedErr}}, nil, expectedErr},
		{"returns existing source company", &Account{Account: account, source: loadedSource}, company, nil},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := test.account.GetCompany(nil)

			assert.Equal(t, test.err, err)
			assert.Equal(t, test.company, result)
		})
	}
}

func assertAccountsError(t *testing.T, result []*Account, err error, expectedErr error) {
	assert.Same(t, expectedErr, err)
	assert.Nil(t, result)
}

func Test_GetAllAccounts(t *testing.T) {
	t.Run("returns accounts", func(t *testing.T) {
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			dbAccounts := []*table.Account{{ID: 1}}
			getAllAccountsStub := mocka.Function(t, &getAllAccounts, dbAccounts, nil)
			defer getAllAccountsStub.Restore()

			result, err := GetAllAccounts(tx)

			assert.Equal(t, []interface{}{tx}, getAllAccountsStub.GetFirstCall().Arguments())
			assert.Nil(t, err)
			assert.Equal(t, dbAccounts[0], result[0].Account)
			assert.NotNil(t, result[0].source)
		})
	})
	t.Run("returns error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			expectedErr := errors.New("database error")
			getAllAccountsStub := mocka.Function(t, &getAllAccounts, nil, expectedErr)
			defer getAllAccountsStub.Restore()

			result, err := GetAllAccounts(tx)

			assertAccountsError(t, result, err, expectedErr)
		})
	})
}

func Test_GetAccountByID(t *testing.T) {
	id := int64(42)
	t.Run("returns accounts", func(t *testing.T) {
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			dbAccounts := []*table.Account{{ID: 1}}
			getAccountByIDStub := mocka.Function(t, &getAccountByID, dbAccounts, nil)
			defer getAccountByIDStub.Restore()

			result, err := GetAccountByID(tx, id)

			assert.Equal(t, []interface{}{tx, id}, getAccountByIDStub.GetFirstCall().Arguments())
			assert.Nil(t, err)
			assert.Equal(t, dbAccounts[0], result[0].Account)
			assert.NotNil(t, result[0].source)
		})
	})
	t.Run("returns error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			expectedErr := errors.New("database error")
			getAccountByIDStub := mocka.Function(t, &getAccountByID, nil, expectedErr)
			defer getAccountByIDStub.Restore()

			result, err := GetAccountByID(tx, id)

			assertAccountsError(t, result, err, expectedErr)
		})
	})
}

func Test_GetAccountsByName(t *testing.T) {
	name := "account name"
	t.Run("returns accounts", func(t *testing.T) {
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			dbAccounts := []*table.Account{{ID: 1}}
			getAccountsByNameStub := mocka.Function(t, &getAccountsByName, dbAccounts, nil)
			defer getAccountsByNameStub.Restore()

			result, err := GetAccountsByName(tx, name)

			assert.Equal(t, []interface{}{tx, name}, getAccountsByNameStub.GetFirstCall().Arguments())
			assert.Nil(t, err)
			assert.Equal(t, dbAccounts[0], result[0].Account)
			assert.NotNil(t, result[0].source)
		})
	})
	t.Run("returns error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			expectedErr := errors.New("database error")
			getAccountsByNameStub := mocka.Function(t, &getAccountsByName, nil, expectedErr)
			defer getAccountsByNameStub.Restore()

			result, err := GetAccountsByName(tx, name)

			assertAccountsError(t, result, err, expectedErr)
		})
	})
}

func Test_GetAccountsByCompanyID(t *testing.T) {
	ids := []int64{42, 69}
	t.Run("returns accounts", func(t *testing.T) {
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			dbAccounts := []*table.Account{{ID: 1}}
			getAccountsByCompanyIDsStub := mocka.Function(t, &getAccountsByCompanyIDs, dbAccounts, nil)
			defer getAccountsByCompanyIDsStub.Restore()

			result, err := GetAccountsByCompanyIDs(tx, ids)

			assert.Equal(t, []interface{}{tx, ids}, getAccountsByCompanyIDsStub.GetFirstCall().Arguments())
			assert.Nil(t, err)
			assert.Equal(t, dbAccounts[0], result[0].Account)
			assert.NotNil(t, result[0].source)
		})
	})
	t.Run("returns error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			expectedErr := errors.New("database error")
			getAccountsByCompanyIDsStub := mocka.Function(t, &getAccountsByCompanyIDs, nil, expectedErr)
			defer getAccountsByCompanyIDsStub.Restore()

			result, err := GetAccountsByCompanyIDs(tx, ids)

			assertAccountsError(t, result, err, expectedErr)
		})
	})
}
