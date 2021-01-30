package model

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

func Test_category_ptrTo(t *testing.T) {
	category := &Category{}
	tests := []struct {
		column string
		ptr    interface{}
	}{
		{column: "id", ptr: &category.ID},
		{column: "code", ptr: &category.Code},
		{column: "description", ptr: &category.Description},
		{column: "amount_type", ptr: &category.AmountType},
		{column: "parent_id", ptr: &category.ParentID},
		{column: "security", ptr: &category.Security},
		{column: "income", ptr: &category.Income},
		{column: "version", ptr: &category.Version},
		{column: "transaction_count", ptr: &category.TransactionCount},
		{column: "change_user", ptr: &category.ChangeUser},
		{column: "change_date", ptr: &category.ChangeDate},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s returns pointer to field", test.column), func(t *testing.T) {
			field := category.ptrTo(test.column)

			assert.Same(t, test.ptr, field)
		})
	}
}

func Test_GetAllCategories(t *testing.T) {
	t.Run("returns categories", func(t *testing.T) {
		categories := []*Category{{ID: 1}}
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
