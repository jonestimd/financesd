package table

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Account_PtrTo(t *testing.T) {
	account := &Account{}
	tests := []struct {
		column string
		ptr    interface{}
	}{
		{column: "id", ptr: &account.ID},
		{column: "company_id", ptr: &account.CompanyID},
		{column: "name", ptr: &account.Name},
		{column: "description", ptr: &account.Description},
		{column: "account_no", ptr: &account.AccountNo},
		{column: "type", ptr: &account.Type},
		{column: "closed", ptr: &account.Closed},
		{column: "currency_id", ptr: &account.CurrencyID},
		{column: "version", ptr: &account.Version},
		{column: "balance", ptr: &account.Balance},
		{column: "transaction_count", ptr: &account.TransactionCount},
		{column: "change_user", ptr: &account.ChangeUser},
		{column: "change_date", ptr: &account.ChangeDate},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s returns pointer to field", test.column), func(t *testing.T) {
			field := account.PtrTo(test.column)

			assert.Same(t, test.ptr, field)
		})
	}
}
