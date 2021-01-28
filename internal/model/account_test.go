package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MonsantoCo/mocka/v2"
	"github.com/jonestimd/financesd/internal/sqltest"
	"github.com/stretchr/testify/assert"
)

func Test_account_ptrTo(t *testing.T) {
	account := &Account{}
	tests := []struct {
		column string
		ptr    interface{}
	}{
		{column: "id", ptr: &account.ID},
		{column: "company_id", ptr: &account.CompanyID},
		{column: "name", ptr: &account.Name},
		{column: "description", ptr: &account.Description},
		{column: "account_no", ptr: &account.AccountNo},
		{column: "type", ptr: &account.Type},
		{column: "closed", ptr: &account.Closed},
		{column: "currency_id", ptr: &account.CurrencyID},
		{column: "version", ptr: &account.Version},
		{column: "balance", ptr: &account.Balance},
		{column: "transaction_count", ptr: &account.TransactionCount},
		{column: "change_user", ptr: &account.ChangeUser},
		{column: "change_date", ptr: &account.ChangeDate},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s returns pointer to field", test.column), func(t *testing.T) {
			field := account.ptrTo(test.column)

			assert.Same(t, test.ptr, field)
		})
	}
}

func Test_Account_GetCompany(t *testing.T) {
	companyID := int64(42)
	company := &Company{ID: companyID}
	loadedSource := &companySource{companiesByID: map[int64]*Company{companyID: company}}
	expectedErr := errors.New("test error")
	tests := []struct {
		name    string
		account *Account
		company *Company
		err     error
	}{
		{"returns nil for nil company ID", &Account{}, nil, nil},
		{"returns existing source error", &Account{CompanyID: &companyID, source: &companySource{err: expectedErr}}, nil, expectedErr},
		{"returns existing source company", &Account{CompanyID: &companyID, source: loadedSource}, company, nil},
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
		accounts := []*Account{{ID: 1}}
		runQueryStub := mocka.Function(t, &runQuery, accounts, nil)
		defer runQueryStub.Restore()
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result, err := GetAllAccounts(tx)

			assert.Equal(t, []interface{}{tx, accountType, accountSQL, []interface{}(nil)}, runQueryStub.GetFirstCall().Arguments())
			assert.Nil(t, err)
			assert.Equal(t, accounts, result)
		})
	})
	t.Run("returns error", func(t *testing.T) {
		expectedErr := errors.New("database error")
		runQueryStub := mocka.Function(t, &runQuery, nil, expectedErr)
		defer runQueryStub.Restore()
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result, err := GetAllAccounts(tx)

			assertAccountsError(t, result, err, expectedErr)
		})
	})
}

func Test_GetAccountByID(t *testing.T) {
	id := int64(42)
	t.Run("returns accounts", func(t *testing.T) {
		accounts := []*Account{{ID: 1}}
		runQueryStub := mocka.Function(t, &runQuery, accounts, nil)
		defer runQueryStub.Restore()
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result, err := GetAccountByID(tx, id)

			assert.Equal(t, []interface{}{tx, accountType, accountSQL + " where a.id = ?", []interface{}{id}},
				runQueryStub.GetFirstCall().Arguments())
			assert.Nil(t, err)
			assert.Equal(t, accounts, result)
		})
	})
	t.Run("returns error", func(t *testing.T) {
		expectedErr := errors.New("database error")
		runQueryStub := mocka.Function(t, &runQuery, nil, expectedErr)
		defer runQueryStub.Restore()
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result, err := GetAccountByID(tx, id)

			assertAccountsError(t, result, err, expectedErr)
		})
	})
}

func Test_GetAccountsByName(t *testing.T) {
	name := "account name"
	t.Run("returns accounts", func(t *testing.T) {
		accounts := []*Account{{ID: 1}}
		runQueryStub := mocka.Function(t, &runQuery, accounts, nil)
		defer runQueryStub.Restore()
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result, err := GetAccountsByName(tx, name)

			assert.Equal(t, []interface{}{tx, accountType, accountSQL + " where a.name = ?", []interface{}{name}},
				runQueryStub.GetFirstCall().Arguments())
			assert.Nil(t, err)
			assert.Equal(t, accounts, result)
		})
	})
	t.Run("returns error", func(t *testing.T) {
		expectedErr := errors.New("database error")
		runQueryStub := mocka.Function(t, &runQuery, nil, expectedErr)
		defer runQueryStub.Restore()
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result, err := GetAccountsByName(tx, name)

			assertAccountsError(t, result, err, expectedErr)
		})
	})
}

func Test_GetAccountsByCompanyID(t *testing.T) {
	ids := []int64{42, 69}
	t.Run("returns accounts", func(t *testing.T) {
		accounts := []*Account{{ID: 1}}
		runQueryStub := mocka.Function(t, &runQuery, accounts, nil)
		defer runQueryStub.Restore()
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result, err := getAccountsByCompanyIDs(tx, ids)

			jsonIDs, _ := json.Marshal(ids)
			assert.Equal(t, []interface{}{tx, accountType, accountSQL + " where json_contains(?, cast(company_id as json))", []interface{}{jsonIDs}},
				runQueryStub.GetFirstCall().Arguments())
			assert.Nil(t, err)
			assert.Equal(t, accounts, result)
		})
	})
	t.Run("returns error", func(t *testing.T) {
		expectedErr := errors.New("database error")
		runQueryStub := mocka.Function(t, &runQuery, nil, expectedErr)
		defer runQueryStub.Restore()
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result, err := getAccountsByCompanyIDs(tx, ids)

			assertAccountsError(t, result, err, expectedErr)
		})
	})
}
