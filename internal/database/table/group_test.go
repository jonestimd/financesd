package table

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Group_PtrTo(t *testing.T) {
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
			field := group.PtrTo(test.column)

			assert.Same(t, test.ptr, field)
		})
	}
}
