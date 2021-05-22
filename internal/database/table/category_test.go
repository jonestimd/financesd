package table

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_category_PtrTo(t *testing.T) {
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
			field := category.PtrTo(test.column)

			assert.Same(t, test.ptr, field)
		})
	}
}
