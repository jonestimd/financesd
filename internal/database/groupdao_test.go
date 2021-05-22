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

func Test_GetAllGroups(t *testing.T) {
	t.Run("returns groups", func(t *testing.T) {
		groups := []*table.Group{{ID: 1}}
		runQueryStub := mocka.Function(t, &runQuery, groups, nil)
		defer runQueryStub.Restore()
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result, err := GetAllGroups(tx)

			assert.Nil(t, err)
			assert.Equal(t, []interface{}{tx, groupType, groupSQL, []interface{}(nil)}, runQueryStub.GetFirstCall().Arguments())
			assert.Equal(t, groups, result)
		})
	})
	t.Run("returns error", func(t *testing.T) {
		expectedErr := errors.New("database error")
		runQueryStub := mocka.Function(t, &runQuery, nil, expectedErr)
		defer runQueryStub.Restore()
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result, err := GetAllGroups(tx)

			assert.Same(t, expectedErr, err)
			assert.Nil(t, result)
		})
	})
}
