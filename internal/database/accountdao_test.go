package database

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MonsantoCo/mocka/v2"
	"github.com/jonestimd/financesd/internal/database/table"
	"github.com/jonestimd/financesd/internal/sqltest"
	"github.com/stretchr/testify/assert"
)

func Test_GetAllAccounts(t *testing.T) {
	t.Run("returns accounts", func(t *testing.T) {
		accounts := []*table.Account{{ID: 1}}
		runQueryStub := mocka.Function(t, &runQuery, accounts)
		defer runQueryStub.Restore()
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result := GetAllAccounts(tx)

			assert.Equal(t, []interface{}{tx, accountType, accountSQL, []interface{}(nil)}, runQueryStub.GetFirstCall().Arguments())
			assert.Equal(t, accounts, result)
		})
	})
}

func Test_GetAccountByID(t *testing.T) {
	id := int64(42)
	accounts := []*table.Account{{ID: 1}}
	runQueryStub := mocka.Function(t, &runQuery, accounts)
	defer runQueryStub.Restore()
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		result := GetAccountByID(tx, id)

		assert.Equal(t, []interface{}{tx, accountType, accountSQL + " where a.id = ?", []interface{}{id}},
			runQueryStub.GetFirstCall().Arguments())
		assert.Equal(t, accounts, result)
	})
}

func Test_GetAccountsByName(t *testing.T) {
	name := "account name"
	accounts := []*table.Account{{ID: 1}}
	runQueryStub := mocka.Function(t, &runQuery, accounts)
	defer runQueryStub.Restore()
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		result := GetAccountsByName(tx, name)

		assert.Equal(t, []interface{}{tx, accountType, accountSQL + " where a.name = ?", []interface{}{name}},
			runQueryStub.GetFirstCall().Arguments())
		assert.Equal(t, accounts, result)
	})
}

func Test_GetAccountsByCompanyID(t *testing.T) {
	ids := []int64{42, 69}
	accounts := []*table.Account{{ID: 1}}
	runQueryStub := mocka.Function(t, &runQuery, accounts)
	defer runQueryStub.Restore()
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		result := GetAccountsByCompanyIDs(tx, ids)

		jsonIDs, _ := json.Marshal(ids)
		assert.Equal(t, []interface{}{tx, accountType, accountSQL + " where json_contains(?, cast(company_id as json))", []interface{}{jsonIDs}},
			runQueryStub.GetFirstCall().Arguments())
		assert.Equal(t, accounts, result)
	})
}
