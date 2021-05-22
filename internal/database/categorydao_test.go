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

func Test_GetAllCategories(t *testing.T) {
	t.Run("returns categories", func(t *testing.T) {
		categories := []*table.Category{{ID: 1}}
		runQueryStub := mocka.Function(t, &runQuery, categories, nil)
		defer runQueryStub.Restore()
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result, err := GetAllCategories(tx)

			assert.Nil(t, err)
			assert.Equal(t, []interface{}{tx, categoryType, categorySQL, []interface{}(nil)}, runQueryStub.GetFirstCall().Arguments())
			assert.Equal(t, categories, result)
		})
	})
	t.Run("returns error", func(t *testing.T) {
		expectedErr := errors.New("database error")
		runQueryStub := mocka.Function(t, &runQuery, nil, expectedErr)
		defer runQueryStub.Restore()
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result, err := GetAllCategories(tx)

			assert.Same(t, expectedErr, err)
			assert.Nil(t, result)
		})
	})
}
