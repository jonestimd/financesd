package graphql

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MonsantoCo/mocka/v2"
	"github.com/jonestimd/financesd/internal/model"
	"github.com/jonestimd/financesd/internal/sqltest"
	"github.com/stretchr/testify/assert"
)

func Test_groupQueryFields_Resolve(t *testing.T) {
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		groups := []*model.Group{{ID: 42}}
		getAll := mocka.Function(t, &getAllGroups, groups, nil)
		defer getAll.Restore()
		params := newResolveParams(tx, groupQuery, newField("", "id"), newField("", "name"))

		result, err := groupQueryFields.Resolve(params.ResolveParams)

		assert.Nil(t, err)
		assert.Equal(t, groups, result)
		assert.Equal(t, []interface{}{tx}, getAll.GetFirstCall().Arguments())
	})
}
