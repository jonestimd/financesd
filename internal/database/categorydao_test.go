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

func Test_GetAllCategories(t *testing.T) {
	categories := []*table.Category{{ID: 1}}
	runQueryStub := mocka.Function(t, &runQuery, categories)
	defer runQueryStub.Restore()
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		result := GetAllCategories(tx)

		assert.Equal(t, []interface{}{tx, categoryType, categorySQL, []interface{}(nil)}, runQueryStub.GetFirstCall().Arguments())
		assert.Equal(t, categories, result)
	})
}
