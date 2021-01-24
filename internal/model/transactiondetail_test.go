package model

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jonestimd/financesd/internal/sqltest"
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
	t.Run("returns nil for nil RelatedDetailID", func(t *testing.T) {
		result, err := (&TransactionDetail{}).GetRelatedDetail(nil)

		assert.Nil(t, err)
		assert.Nil(t, result)
	})
	t.Run("returns existing error", func(t *testing.T) {
		accountTransactions := &AccountTransactions{err: errors.New("test error")}

		result, err := (&TransactionDetail{RelatedDetailID: &id, accountTransactions: accountTransactions}).GetRelatedDetail(nil)

		assert.Same(t, accountTransactions.err, err)
		assert.Nil(t, result)
	})
	t.Run("loads related details", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			accountTx := &AccountTransactions{accountID: 69}
			mockDB.ExpectQuery(relatedDetailsSQL).WithArgs(accountTx.accountID).WillReturnRows(sqltest.MockRows("id").AddRow(id))

			result, err := (&TransactionDetail{RelatedDetailID: &id, accountTransactions: accountTx}).GetRelatedDetail(tx)

			assert.Nil(t, err)
			assert.Equal(t, &TransactionDetail{ID: id, accountTransactions: accountTx}, result)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})
	})
	t.Run("returns existing detail", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			relatedDetail := &TransactionDetail{ID: id}
			accountTx := &AccountTransactions{accountID: 69, relatedDetailsByID: map[int64]*TransactionDetail{id: relatedDetail}}

			result, err := (&TransactionDetail{RelatedDetailID: &id, accountTransactions: accountTx}).GetRelatedDetail(tx)

			assert.Nil(t, err)
			assert.Same(t, relatedDetail, result)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})
	})
	t.Run("returns query error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			expectedErr := errors.New("test error")
			accountTx := &AccountTransactions{accountID: 69}
			mockDB.ExpectQuery(relatedDetailsSQL).WithArgs(accountTx.accountID).WillReturnError(expectedErr)

			result, err := (&TransactionDetail{RelatedDetailID: &id, accountTransactions: accountTx}).GetRelatedDetail(tx)

			assert.Same(t, expectedErr, err)
			assert.Nil(t, result)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})
	})
}

func Test_GetRelatedTransaction(t *testing.T) {
	id := int64(42)
	t.Run("returns existing error", func(t *testing.T) {
		accountTransactions := &AccountTransactions{err: errors.New("test error")}

		result, err := (&TransactionDetail{RelatedDetailID: &id, accountTransactions: accountTransactions}).GetRelatedTransaction(nil)

		assert.Same(t, accountTransactions.err, err)
		assert.Nil(t, result)
	})
	t.Run("loads related transactions", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			accountTx := &AccountTransactions{accountID: 69}
			mockDB.ExpectQuery(relatedTransactionSQL).WithArgs(accountTx.accountID).WillReturnRows(sqltest.MockRows("id").AddRow(id))

			result, err := (&TransactionDetail{TransactionID: id, accountTransactions: accountTx}).GetRelatedTransaction(tx)

			assert.Nil(t, err)
			assert.Equal(t, &Transaction{ID: id}, result)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})
	})
	t.Run("returns existing transaction", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			relatedTx := &Transaction{ID: id}
			accountTx := &AccountTransactions{accountID: 69, relatedTxByID: map[int64]*Transaction{id: relatedTx}}

			result, err := (&TransactionDetail{TransactionID: id, accountTransactions: accountTx}).GetRelatedTransaction(tx)

			assert.Nil(t, err)
			assert.Same(t, relatedTx, result)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})
	})
	t.Run("returns query error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			expectedErr := errors.New("test error")
			accountTx := &AccountTransactions{accountID: 69}
			mockDB.ExpectQuery(relatedTransactionSQL).WithArgs(accountTx.accountID).WillReturnError(expectedErr)

			result, err := (&TransactionDetail{TransactionID: id, accountTransactions: accountTx}).GetRelatedTransaction(tx)

			assert.Same(t, expectedErr, err)
			assert.Nil(t, result)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})
	})
}
