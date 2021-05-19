package database

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_TransactionDetail_ptrTo(t *testing.T) {
	detail := &TransactionDetail{}
	tests := []struct {
		column string
		ptr    interface{}
	}{
		{column: "id", ptr: &detail.ID},
		{column: "transaction_id", ptr: &detail.TransactionID},
		{column: "transaction_category_id", ptr: &detail.TransactionCategoryID},
		{column: "transaction_group_id", ptr: &detail.TransactionGroupID},
		{column: "memo", ptr: &detail.Memo},
		{column: "amount", ptr: &detail.Amount},
		{column: "asset_quantity", ptr: &detail.AssetQuantity},
		{column: "exchange_asset_id", ptr: &detail.ExchangeAssetID},
		{column: "related_detail_id", ptr: &detail.RelatedDetailID},
		{column: "version", ptr: &detail.Version},
		{column: "change_user", ptr: &detail.ChangeUser},
		{column: "change_date", ptr: &detail.ChangeDate},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s returns pointer to field", test.column), func(t *testing.T) {
			field := detail.ptrTo(test.column)

			assert.Same(t, test.ptr, field)
		})
	}
}
