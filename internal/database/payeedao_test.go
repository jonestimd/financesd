package database

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

func Test_GetAllPayees(t *testing.T) {
	t.Run("returns payees", func(t *testing.T) {
		payees := []*table.Payee{{ID: 1}}
		runQueryStub := mocka.Function(t, &runQuery, payees, nil)
		defer runQueryStub.Restore()
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result, err := GetAllPayees(tx)

			assert.Nil(t, err)
			assert.Equal(t, []interface{}{tx, payeeType, payeeSQL, []interface{}(nil)}, runQueryStub.GetFirstCall().Arguments())
			assert.Equal(t, payees, result)
		})
	})
	t.Run("returns error", func(t *testing.T) {
		expectedErr := errors.New("database error")
		runQueryStub := mocka.Function(t, &runQuery, nil, expectedErr)
		defer runQueryStub.Restore()
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result, err := GetAllPayees(tx)

			assert.Same(t, expectedErr, err)
			assert.Nil(t, result)
		})
	})
}
