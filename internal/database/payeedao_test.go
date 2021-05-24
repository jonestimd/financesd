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

func Test_GetAllPayees(t *testing.T) {
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		payees := []*table.Payee{{ID: 1}}
		runQueryStub := mocka.Function(t, &runQuery, payees)
		defer runQueryStub.Restore()

		result := GetAllPayees(tx)

		assert.Equal(t, []interface{}{tx, payeeType, payeeSQL, []interface{}(nil)}, runQueryStub.GetFirstCall().Arguments())
		assert.Equal(t, payees, result)
	})
}
