package model

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MonsantoCo/mocka/v2"
	"github.com/jonestimd/financesd/internal/sqltest"
	"github.com/stretchr/testify/assert"
)

func Test_security_ptrTo(t *testing.T) {
	security := &Security{}
	tests := []struct {
		column string
		ptr    interface{}
	}{
		{column: "asset_id", ptr: &security.AssetID},
		{column: "security_type", ptr: &security.Type},
		{column: "id", ptr: &security.ID},
		{column: "name", ptr: &security.Name},
		{column: "type", ptr: &security.Asset.Type},
		{column: "shares", ptr: &security.Shares},
		{column: "first_acquired", ptr: &security.FirstAcquired},
		{column: "cost_basis", ptr: &security.CostBasis},
		{column: "dividends", ptr: &security.Dividends},
		{column: "transaction_count", ptr: &security.TransactionCount},
		{column: "scale", ptr: &security.Scale},
		{column: "symbol", ptr: &security.Symbol},
		{column: "version", ptr: &security.Version},
		{column: "change_user", ptr: &security.ChangeUser},
		{column: "change_date", ptr: &security.ChangeDate},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s returns pointer to field", test.column), func(t *testing.T) {
			field := security.ptrTo(test.column)

			assert.Same(t, test.ptr, field)
		})
	}
}

func Test_GetAllSecurities(t *testing.T) {
	t.Run("returns securities", func(t *testing.T) {
		securities := []*Security{{AssetID: 1}}
		runQueryStub := mocka.Function(t, &runQuery, securities, nil)
		defer runQueryStub.Restore()
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result, err := GetAllSecurities(tx)

			assert.Nil(t, err)
			assert.Equal(t, []interface{}{tx, securityType, securitySQL, []interface{}(nil)}, runQueryStub.GetFirstCall().Arguments())
			assert.Equal(t, securities, result)
		})
	})
	t.Run("returns error", func(t *testing.T) {
		expectedErr := errors.New("database error")
		runQueryStub := mocka.Function(t, &runQuery, nil, expectedErr)
		defer runQueryStub.Restore()
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result, err := GetAllSecurities(tx)

			assert.Same(t, expectedErr, err)
			assert.Nil(t, result)
		})
	})
}

func Test_GetSecurityByID(t *testing.T) {
	id := int64(42)
	t.Run("returns security", func(t *testing.T) {
		securities := []*Security{{AssetID: id}}
		runQueryStub := mocka.Function(t, &runQuery, securities, nil)
		defer runQueryStub.Restore()
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result, err := GetSecurityByID(tx, id)

			assert.Nil(t, err)
			assert.Equal(t, []interface{}{tx, securityType, securitySQL + " where a.id = ?",
				[]interface{}{id}}, runQueryStub.GetFirstCall().Arguments())
			assert.Equal(t, securities, result)
		})
	})
	t.Run("returns error", func(t *testing.T) {
		expectedErr := errors.New("database error")
		runQueryStub := mocka.Function(t, &runQuery, nil, expectedErr)
		defer runQueryStub.Restore()
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result, err := GetSecurityByID(tx, id)

			assert.Same(t, expectedErr, err)
			assert.Nil(t, result)
		})
	})
}

func Test_GetSecurityBySymbol(t *testing.T) {
	symbol := "S1"
	t.Run("returns security", func(t *testing.T) {
		securities := []*Security{{AssetID: 42}}
		runQueryStub := mocka.Function(t, &runQuery, securities, nil)
		defer runQueryStub.Restore()
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result, err := GetSecurityBySymbol(tx, symbol)

			assert.Nil(t, err)
			assert.Equal(t, []interface{}{tx, securityType, securitySQL + " where s.symbol = ?",
				[]interface{}{symbol}}, runQueryStub.GetFirstCall().Arguments())
			assert.Equal(t, securities, result)
		})
	})
	t.Run("returns error", func(t *testing.T) {
		expectedErr := errors.New("database error")
		runQueryStub := mocka.Function(t, &runQuery, nil, expectedErr)
		defer runQueryStub.Restore()
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result, err := GetSecurityBySymbol(tx, symbol)

			assert.Same(t, expectedErr, err)
			assert.Nil(t, result)
		})
	})
}
