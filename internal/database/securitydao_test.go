package database

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MonsantoCo/mocka/v2"
	"github.com/jonestimd/financesd/internal/database/table"
	"github.com/jonestimd/financesd/internal/sqltest"
	"github.com/stretchr/testify/assert"
)

func Test_GetAllSecurities(t *testing.T) {
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		securities := []*table.Security{{AssetID: 1}}
		runQueryStub := mocka.Function(t, &runQuery, securities)
		defer runQueryStub.Restore()

		result := GetAllSecurities(tx)

		assert.Equal(t, []interface{}{tx, securityType, securitySQL, []interface{}(nil)}, runQueryStub.GetFirstCall().Arguments())
		assert.Equal(t, securities, result)
	})
}

func Test_GetSecurityByID(t *testing.T) {
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		id := int64(42)
		securities := []*table.Security{{AssetID: id}}
		runQueryStub := mocka.Function(t, &runQuery, securities)
		defer runQueryStub.Restore()

		result := GetSecurityByID(tx, id)

		assert.Equal(t, []interface{}{tx, securityType, securitySQL + " where a.id = ?",
			[]interface{}{id}}, runQueryStub.GetFirstCall().Arguments())
		assert.Equal(t, securities, result)
	})
}

func Test_GetSecurityBySymbol(t *testing.T) {
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		symbol := "S1"
		securities := []*table.Security{{AssetID: 42}}
		runQueryStub := mocka.Function(t, &runQuery, securities)
		defer runQueryStub.Restore()

		result := GetSecurityBySymbol(tx, symbol)

		assert.Equal(t, []interface{}{tx, securityType, securitySQL + " where s.symbol = ?",
			[]interface{}{symbol}}, runQueryStub.GetFirstCall().Arguments())
		assert.Equal(t, securities, result)
	})
}
