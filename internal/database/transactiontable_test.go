package database

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Transaction_ptrTo(t *testing.T) {
	transaction := &Transaction{}
	tests := []struct {
		column string
		ptr    interface{}
	}{
		{column: "id", ptr: &transaction.ID},
		{column: "date", ptr: &transaction.Date},
		{column: "memo", ptr: &transaction.Memo},
		{column: "reference_number", ptr: &transaction.ReferenceNumber},
		{column: "cleared", ptr: &transaction.Cleared},
		{column: "account_id", ptr: &transaction.AccountID},
		{column: "payee_id", ptr: &transaction.PayeeID},
		{column: "security_id", ptr: &transaction.SecurityID},
		{column: "version", ptr: &transaction.Version},
		{column: "change_user", ptr: &transaction.ChangeUser},
		{column: "change_date", ptr: &transaction.ChangeDate},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s returns pointer to field", test.column), func(t *testing.T) {
			field := transaction.ptrTo(test.column)

			assert.Same(t, test.ptr, field)
		})
	}
}
