package database

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

func Test_Group_ptrTo(t *testing.T) {
	group := &Group{}
	tests := []struct {
		column string
		ptr    interface{}
	}{
		{column: "id", ptr: &group.ID},
		{column: "name", ptr: &group.Name},
		{column: "description", ptr: &group.Description},
		{column: "version", ptr: &group.Version},
		{column: "transaction_count", ptr: &group.TransactionCount},
		{column: "change_user", ptr: &group.ChangeUser},
		{column: "change_date", ptr: &group.ChangeDate},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s returns pointer to field", test.column), func(t *testing.T) {
			field := group.ptrTo(test.column)

			assert.Same(t, test.ptr, field)
		})
	}
}

func Test_GetAllGroups(t *testing.T) {
	t.Run("returns groups", func(t *testing.T) {
		groups := []*Group{{ID: 1}}
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
