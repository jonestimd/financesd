package table

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Payee_PtrTo(t *testing.T) {
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
			field := payee.PtrTo(test.column)

			assert.Same(t, test.ptr, field)
		})
	}
}
