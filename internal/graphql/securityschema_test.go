package graphql

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jonestimd/financesd/internal/model"
)

func Test_securityQueryFields_Resolve_all(t *testing.T) {
	testQuery(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		var assetID int64 = 1
		params := newResolveParams(tx, securityQuery, nil, newField("", "type"), newField("", "id"), newField("", "name"))
		mock.ExpectQuery(securitySQL).WithArgs().WillReturnRows(
			mockRows("id", "type", "security_type").AddRow(assetID, "Security", "Stock"))
		expectedRows := []map[string]interface{}{
			{"id": assetID, "assetType": "Security", "type": "Stock"},
		}

		result, err := securityQueryFields.Resolve(params)

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

func Test_securityQueryFields_Resolve_queryError(t *testing.T) {
	testQuery(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		params := newResolveParams(tx, securityQuery, nil, newField("", "type"), newField("", "id"), newField("", "name"))
		queryError := errors.New("invalid SQL")
		mock.ExpectQuery(securitySQL).WithArgs().WillReturnError(queryError)

		_, err := securityQueryFields.Resolve(params)

		if err != queryError {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}

func Test_securityQueryFields_Resolve_byID(t *testing.T) {
	testQuery(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		var assetID int64 = 1
		securityType := "Stock"
		name := "Security 1"
		args := map[string]interface{}{"id": fmt.Sprint(assetID)}
		params := newResolveParams(tx, securityQuery, args, newField("", "type"), newField("", "id"), newField("", "name"))
		mock.ExpectQuery(securitySQL + " where a.id = ?").
			WithArgs(assetID).
			WillReturnRows(mockRows("id", "security_type", "name").AddRow(assetID, securityType, name))
		expectedRows := []map[string]interface{}{
			{"id": assetID, "type": securityType, "name": name},
		}

		result, err := securityQueryFields.Resolve(params)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !reflect.DeepEqual(result, expectedRows) {
			t.Errorf("Expected result: %v, got %v", expectedRows[0], result.([]*model.Security)[0])
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
	testQuery(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		args := map[string]interface{}{"id": "one"}
		params := newResolveParams(tx, securityQuery, args, newField("", "type"), newField("", "id"), newField("", "name"))

		_, err := securityQueryFields.Resolve(params)

		if err == nil {
			t.Error("Expected a parsing error")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func Test_securityQueryFields_Resolve_bySymbol(t *testing.T) {
	testQuery(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		var assetID int64 = 1
		symbol := "A1"
		securityType := "Stock"
		name := "Security 1"
		args := map[string]interface{}{"symbol": symbol}
		params := newResolveParams(tx, securityQuery, args, newField("", "type"), newField("", "id"), newField("", "name"))
		mock.ExpectQuery(securitySQL + " where a.symbol = ?").
			WithArgs(symbol).
			WillReturnRows(mockRows("id", "security_type", "name").AddRow(assetID, securityType, name))
		expectedRows := []map[string]interface{}{
			{"id": assetID, "type": securityType, "name": name},
		}

		result, err := securityQueryFields.Resolve(params)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !reflect.DeepEqual(result, expectedRows) {
			t.Errorf("Expected result: %v, got %v", expectedRows[0], result.([]*model.Security)[0])
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
