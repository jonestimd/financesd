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

func Test_Payee_ptrTo(t *testing.T) {
	payee := &Payee{}
	tests := []struct {
		column string
		ptr    interface{}
	}{
		{column: "id", ptr: &payee.ID},
		{column: "name", ptr: &payee.Name},
		{column: "version", ptr: &payee.Version},
		{column: "transaction_count", ptr: &payee.TransactionCount},
		{column: "change_user", ptr: &payee.ChangeUser},
		{column: "change_date", ptr: &payee.ChangeDate},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s returns pointer to field", test.column), func(t *testing.T) {
			field := payee.ptrTo(test.column)

			assert.Same(t, test.ptr, field)
		})
	}
}

func Test_GetAllPayees(t *testing.T) {
	t.Run("returns payees", func(t *testing.T) {
		payees := []*Payee{{ID: 1}}
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
