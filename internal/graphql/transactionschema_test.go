package graphql

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func Test_transactionQueryFields_Resolve_returnsRows(t *testing.T) {
	testQuery(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		accountID := "123"
		args := map[string]interface{}{"accountId": accountID}
		params := newResolveParams(tx, transactionQuery, args, newField("", "id"), newField("", "memo"))
		rows := sqlmock.NewRows([]string{"json"}).AddRow(`{"id":1,"name":"Transaction 1"}`)
		mock.ExpectQuery(`select json_object("id", t.id, "memo", t.memo) from transaction t where t.account_id = ? order by t.date, t.id`).
			WithArgs(accountID).
			WillReturnRows(rows)
		expectedRows := []map[string]interface{}{
			{"id": 1.0, "name": "Transaction 1"},
		}

		result, err := transactionQueryFields.Resolve(params)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !reflect.DeepEqual(result, expectedRows) {
			t.Errorf("Expected result: %v, got %v", expectedRows, result)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func Test_transactionQueryFields_Resolve_requiresAccountID(t *testing.T) {
	testQuery(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		params := newResolveParams(tx, transactionQuery, nil, newField("", "id"), newField("", "memo"))

		_, err := transactionQueryFields.Resolve(params)

		if err == nil {
			t.Error("Expected an error")
		}
		if fmt.Sprint(err) != "accountId is required" {
			t.Errorf("Unexpected message: %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
