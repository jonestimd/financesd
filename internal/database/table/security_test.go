package table

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_security_PtrTo(t *testing.T) {
	security := &Security{}
	tests := []struct {
		column string
		ptr    interface{}
	}{
		{column: "asset_id", ptr: &security.AssetID},
		{column: "security_type", ptr: &security.Type},
		{column: "id", ptr: &security.ID},
		{column: "name", ptr: &security.Name},
		{column: "type", ptr: &security.Asset.Type},
		{column: "shares", ptr: &security.Shares},
		{column: "first_acquired", ptr: &security.FirstAcquired},
		{column: "cost_basis", ptr: &security.CostBasis},
		{column: "dividends", ptr: &security.Dividends},
		{column: "transaction_count", ptr: &security.TransactionCount},
		{column: "scale", ptr: &security.Scale},
		{column: "symbol", ptr: &security.Symbol},
		{column: "version", ptr: &security.Version},
		{column: "change_user", ptr: &security.ChangeUser},
		{column: "change_date", ptr: &security.ChangeDate},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s returns pointer to field", test.column), func(t *testing.T) {
			field := security.PtrTo(test.column)

			assert.Same(t, test.ptr, field)
		})
	}
}
