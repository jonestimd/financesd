package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MonsantoCo/mocka/v2"
	"github.com/jonestimd/financesd/internal/database/table"
	"github.com/jonestimd/financesd/internal/sqltest"
	"github.com/stretchr/testify/assert"
)

func assertAccountsError(t *testing.T, result []*table.Account, err error, expectedErr error) {
	assert.Same(t, expectedErr, err)
	assert.Nil(t, result)
}

func Test_GetAllAccounts(t *testing.T) {
	t.Run("returns accounts", func(t *testing.T) {
		accounts := []*table.Account{{ID: 1}}
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
		accounts := []*table.Account{{ID: 1}}
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
		accounts := []*table.Account{{ID: 1}}
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
		accounts := []*table.Account{{ID: 1}}
		runQueryStub := mocka.Function(t, &runQuery, accounts, nil)
		defer runQueryStub.Restore()
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result, err := GetAccountsByCompanyIDs(tx, ids)

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
			result, err := GetAccountsByCompanyIDs(tx, ids)

			assertAccountsError(t, result, err, expectedErr)
		})
	})
}
