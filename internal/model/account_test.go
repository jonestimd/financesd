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

func Test_Accounts(t *testing.T) {
	companyID1 := int64(10)
	companyID2 := int64(20)
	accounts := []*Account{
		{ID: 1, CompanyID: &companyID1},
		{ID: 2, CompanyID: &companyID1},
		{ID: 3},
		{ID: 4, CompanyID: &companyID2},
	}
	a := &Accounts{Accounts: accounts}
	t.Run("CompanyIDs populates byID", func(t *testing.T) {
		companyIDs := a.CompanyIDs()

		assert.ElementsMatch(t, []int64{companyID1, companyID2}, companyIDs)
		assert.Equal(t, a.byID, map[int64]*Account{
			1: accounts[0], 2: accounts[1], 3: accounts[2], 4: accounts[3],
		})
	})
	t.Run("ByID populates companyIDs", func(t *testing.T) {
		a.byID = nil
		a.companyIDs = nil

		byID := a.ByID()

		assert.Equal(t, byID, map[int64]*Account{
			1: accounts[0], 2: accounts[1], 3: accounts[2], 4: accounts[3],
		})
		assert.ElementsMatch(t, []int64{companyID1, companyID2}, a.companyIDs)
	})
	t.Run("Result returns accounts", func(t *testing.T) {
		assert.Equal(t, a.Result(), accounts)
	})
	t.Run("Error returns err", func(t *testing.T) {
		err := errors.New("test error")

		assert.Nil(t, a.Error())
		assert.Same(t, NewAccounts(err).Error(), err)
	})
}

func assertAccountsError(t *testing.T, result *Accounts, err error) {
	assert.Same(t, err, result.Error())
	assert.Nil(t, result.Accounts)
}

func Test_GetAllAccounts(t *testing.T) {
	t.Run("returns accounts", func(t *testing.T) {
		accounts := []*Account{{ID: 1}}
		runQueryStub := mocka.Function(t, &runQuery, accounts, nil)
		defer runQueryStub.Restore()
		sqltest.TestQuery(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result := GetAllAccounts(tx)

			assert.Equal(t, []interface{}{tx, accountType, accountsSQL, []interface{}(nil)}, runQueryStub.GetFirstCall().Arguments())
			assert.Nil(t, result.Error())
			assert.Equal(t, accounts, result.Accounts)
		})
	})
	t.Run("returns error", func(t *testing.T) {
		err := errors.New("database error")
		runQueryStub := mocka.Function(t, &runQuery, nil, err)
		defer runQueryStub.Restore()
		sqltest.TestQuery(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result := GetAllAccounts(tx)

			assertAccountsError(t, result, err)
		})
	})
}

func Test_GetAccountByID(t *testing.T) {
	id := int64(42)
	t.Run("returns accounts", func(t *testing.T) {
		accounts := []*Account{{ID: 1}}
		runQueryStub := mocka.Function(t, &runQuery, accounts, nil)
		defer runQueryStub.Restore()
		sqltest.TestQuery(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result := GetAccountByID(tx, id)

			assert.Equal(t, []interface{}{tx, accountType, accountsSQL + " where a.id = ?", []interface{}{id}},
				runQueryStub.GetFirstCall().Arguments())
			assert.Nil(t, result.Error())
			assert.Equal(t, accounts, result.Accounts)
		})
	})
	t.Run("returns error", func(t *testing.T) {
		err := errors.New("database error")
		runQueryStub := mocka.Function(t, &runQuery, nil, err)
		defer runQueryStub.Restore()
		sqltest.TestQuery(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result := GetAccountByID(tx, id)

			assertAccountsError(t, result, err)
		})
	})
}

func Test_GetAccountsByName(t *testing.T) {
	name := "account name"
	t.Run("returns accounts", func(t *testing.T) {
		accounts := []*Account{{ID: 1}}
		runQueryStub := mocka.Function(t, &runQuery, accounts, nil)
		defer runQueryStub.Restore()
		sqltest.TestQuery(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result := GetAccountsByName(tx, name)

			assert.Equal(t, []interface{}{tx, accountType, accountsSQL + " where a.name = ?", []interface{}{name}},
				runQueryStub.GetFirstCall().Arguments())
			assert.Nil(t, result.Error())
			assert.Equal(t, accounts, result.Accounts)
		})
	})
	t.Run("returns error", func(t *testing.T) {
		err := errors.New("database error")
		runQueryStub := mocka.Function(t, &runQuery, nil, err)
		defer runQueryStub.Restore()
		sqltest.TestQuery(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result := GetAccountsByName(tx, name)

			assertAccountsError(t, result, err)
		})
	})
}

func Test_GetAccountsByCompanyID(t *testing.T) {
	ids := []int64{42, 69}
	t.Run("returns accounts", func(t *testing.T) {
		accounts := []*Account{{ID: 1}}
		runQueryStub := mocka.Function(t, &runQuery, accounts, nil)
		defer runQueryStub.Restore()
		sqltest.TestQuery(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result := GetAccountsByCompanyIDs(tx, ids)

			jsonIDs, _ := json.Marshal(ids)
			assert.Equal(t, []interface{}{tx, accountType, accountsSQL + " where json_contains(?, cast(company_id as json))", []interface{}{jsonIDs}},
				runQueryStub.GetFirstCall().Arguments())
			assert.Nil(t, result.Error())
			assert.Equal(t, accounts, result.Accounts)
		})
	})
	t.Run("returns error", func(t *testing.T) {
		err := errors.New("database error")
		runQueryStub := mocka.Function(t, &runQuery, nil, err)
		defer runQueryStub.Restore()
		sqltest.TestQuery(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result := GetAccountsByCompanyIDs(tx, ids)

			assertAccountsError(t, result, err)
		})
	})
}
