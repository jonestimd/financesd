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

func Test_GetAllGroups(t *testing.T) {
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		groups := []*table.Group{{ID: 1}}
		runQueryStub := mocka.Function(t, &runQuery, groups)
		defer runQueryStub.Restore()

		result := GetAllGroups(tx)

		assert.Equal(t, []interface{}{tx, groupType, groupSQL, []interface{}(nil)}, runQueryStub.GetFirstCall().Arguments())
		assert.Equal(t, groups, result)
	})
}
