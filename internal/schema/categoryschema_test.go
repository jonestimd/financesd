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

func Test_categoryQueryFields_Resolve(t *testing.T) {
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		categories := []*table.Category{{ID: 42}}
		getAll := mocka.Function(t, &getAllCategories, categories)
		defer getAll.Restore()
		params := newResolveParams(tx, categoryQuery, newField("", "id"), newField("", "code"))

		result, err := categoryQueryFields.Resolve(params.ResolveParams)

		assert.Nil(t, err)
		assert.Equal(t, categories, result)
		assert.Equal(t, []interface{}{tx}, getAll.GetFirstCall().Arguments())
	})
}
