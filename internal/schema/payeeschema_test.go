package schema

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MonsantoCo/mocka/v2"
	"github.com/jonestimd/financesd/internal/database/table"
	"github.com/jonestimd/financesd/internal/sqltest"
	"github.com/stretchr/testify/assert"
)

func Test_payeeQueryFields_Resolve(t *testing.T) {
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		payees := []*table.Payee{{ID: 1}}
		getAll := mocka.Function(t, &getAllPayees, payees, nil)
		defer getAll.Restore()
		params := newResolveParams(tx, payeeQuery, newField("", "id"), newField("", "name"))

		result, err := payeeQueryFields.Resolve(params.ResolveParams)

		assert.Nil(t, err)
		assert.Equal(t, payees, result)
		assert.Equal(t, []interface{}{tx}, getAll.GetFirstCall().Arguments())
	})
}
