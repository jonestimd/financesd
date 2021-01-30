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

func Test_GetTransactions(t *testing.T) {
	accountID := int64(42)
	t.Run("returns query error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			expectedErr := errors.New("test error")
			mockDB.ExpectQuery(transactionSQL).WithArgs(accountID).WillReturnError(expectedErr)

			_, err := GetTransactions(tx, accountID)

			assert.Same(t, expectedErr, err)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})
	})
	t.Run("returns transactions", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			txID := int64(69)
			expectedTx := &Transaction{
				ID:                  txID,
				accountTransactions: &AccountTransactions{accountID: accountID},
			}
			mockDB.ExpectQuery(transactionSQL).WithArgs(accountID).WillReturnRows(sqltest.MockRows("id").AddRow(txID))

			result, err := GetTransactions(tx, accountID)

			assert.Nil(t, err)
			assert.Equal(t, []*Transaction{expectedTx}, result)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})
	})
}

func Test_GetDetails(t *testing.T) {
	accountID := int64(42)
	txID := int64(96)
	detailID := int64(69)
	t.Run("returns existing error", func(t *testing.T) {
		accountTx := &AccountTransactions{err: errors.New("test error")}
		transaction := &Transaction{accountTransactions: accountTx}

		_, err := transaction.GetDetails(nil)

		assert.Same(t, accountTx.err, err)
	})
	t.Run("returns existing details", func(t *testing.T) {
		detail := &TransactionDetail{ID: detailID, TransactionID: txID}
		accountTx := &AccountTransactions{accountID: accountID, detailsByTxID: map[int64][]*TransactionDetail{txID: {detail}}}
		transaction := &Transaction{ID: txID, accountTransactions: accountTx}

		result, err := transaction.GetDetails(nil)

		assert.Nil(t, err)
		assert.Equal(t, []*TransactionDetail{detail}, result)
	})
	t.Run("loads transaction details", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			accountTx := &AccountTransactions{accountID: accountID}
			transaction := &Transaction{ID: txID, accountTransactions: accountTx}
			expectedDetail := &TransactionDetail{
				ID:                  detailID,
				TransactionID:       txID,
				accountTransactions: accountTx,
			}
			mockDB.ExpectQuery(transactionDetailSQL).WithArgs(accountID).WillReturnRows(sqltest.MockRows("id", "transaction_id").AddRow(detailID, txID))

			result, err := transaction.GetDetails(tx)

			assert.Nil(t, err)
			assert.Equal(t, []*TransactionDetail{expectedDetail}, result)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})
	})
	t.Run("returns query error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			expectedErr := errors.New("test error")
			accountTx := &AccountTransactions{accountID: accountID}
			transaction := &Transaction{ID: txID, accountTransactions: accountTx}
			mockDB.ExpectQuery(transactionDetailSQL).WithArgs(accountID).WillReturnError(expectedErr)

			_, err := transaction.GetDetails(tx)

			assert.Same(t, expectedErr, err)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})
	})
}
