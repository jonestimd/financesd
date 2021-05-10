package model

import (
	"errors"
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

func Test_GetRelatedDetail(t *testing.T) {
	id := int64(42)
	expectedErr := errors.New("test error")
	relatedDetail := &TransactionDetail{ID: id}
	relatedDetailsByID := map[int64]*TransactionDetail{id: relatedDetail}
	tests := []struct {
		name           string
		relatedID      *int64
		txSource       *transactionSource
		expectedErr    error
		expectedResult *TransactionDetail
	}{
		{"returns nil for nil RelatedDetailID", nil, nil, nil, nil},
		{"returns existing error", &id, &transactionSource{err: expectedErr}, expectedErr, nil},
		{"returns existing detail", &id, &transactionSource{relatedDetailsByID: relatedDetailsByID}, nil, relatedDetail},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			detail := &TransactionDetail{RelatedDetailID: test.relatedID, txSource: test.txSource}

			result, err := detail.GetRelatedDetail(nil)

			assert.Equal(t, test.expectedErr, err)
			assert.Equal(t, test.expectedResult, result)
		})
	}
}

func Test_GetRelatedTransaction(t *testing.T) {
	id := int64(42)
	expectedErr := errors.New("test error")
	relatedTx := &Transaction{ID: id}
	relatedTxByID := map[int64]*Transaction{id: relatedTx}
	tests := []struct {
		name           string
		txSource       *transactionSource
		expectedErr    error
		expectedResult *Transaction
	}{
		{"returns existing error", &transactionSource{err: expectedErr}, expectedErr, nil},
		{"returns existing detail", &transactionSource{relatedTxByID: relatedTxByID}, nil, relatedTx},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			detail := &TransactionDetail{TransactionID: id, txSource: test.txSource}

			result, err := detail.GetRelatedTransaction(nil)

			assert.Equal(t, test.expectedErr, err)
			assert.Equal(t, test.expectedResult, result)
		})
	}
}
