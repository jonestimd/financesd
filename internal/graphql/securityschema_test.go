package graphql

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jonestimd/financesd/internal/model"
	"github.com/jonestimd/financesd/internal/sqltest"
)

func Test_securityQueryFields_Resolve_all(t *testing.T) {
	sqltest.TestQuery(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		var assetID int64 = 1
		params := newResolveParams(tx, securityQuery, newField("", "type"), newField("", "id"), newField("", "name"))
		mock.ExpectQuery(securitySQL).WithArgs().WillReturnRows(
			sqltest.MockRows("id", "type", "security_type").AddRow(assetID, "Security", "Stock"))
		expectedRows := []map[string]interface{}{
			{"id": assetID, "assetType": "Security", "type": "Stock"},
		}

		result, err := securityQueryFields.Resolve(params.ResolveParams)

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
	sqltest.TestQuery(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		params := newResolveParams(tx, securityQuery, newField("", "type"), newField("", "id"), newField("", "name"))
		queryError := errors.New("invalid SQL")
		mock.ExpectQuery(securitySQL).WithArgs().WillReturnError(queryError)

		_, err := securityQueryFields.Resolve(params.ResolveParams)

		if err != queryError {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}

func Test_securityQueryFields_Resolve_byID(t *testing.T) {
	sqltest.TestQuery(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		var assetID int64 = 1
		securityType := "Stock"
		name := "Security 1"
		params := newResolveParams(tx, securityQuery, newField("", "type"), newField("", "id"), newField("", "name")).
			addArg("id", fmt.Sprint(assetID))
		mock.ExpectQuery(securitySQL + " where a.id = ?").
			WithArgs(assetID).
			WillReturnRows(sqltest.MockRows("id", "security_type", "name").AddRow(assetID, securityType, name))
		expectedRows := []map[string]interface{}{
			{"id": assetID, "type": securityType, "name": name},
		}

		result, err := securityQueryFields.Resolve(params.ResolveParams)

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
	sqltest.TestQuery(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		params := newResolveParams(tx, securityQuery, newField("", "type"), newField("", "id"), newField("", "name")).addArg("id", "one")

		_, err := securityQueryFields.Resolve(params.ResolveParams)

		if err == nil {
			t.Error("Expected a parsing error")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func Test_securityQueryFields_Resolve_bySymbol(t *testing.T) {
	sqltest.TestQuery(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		var assetID int64 = 1
		symbol := "A1"
		securityType := "Stock"
		name := "Security 1"
		params := newResolveParams(tx, securityQuery, newField("", "type"), newField("", "id"), newField("", "name")).addArg("symbol", symbol)
		mock.ExpectQuery(securitySQL + " where a.symbol = ?").
			WithArgs(symbol).
			WillReturnRows(sqltest.MockRows("id", "security_type", "name").AddRow(assetID, securityType, name))
		expectedRows := []map[string]interface{}{
			{"id": assetID, "type": securityType, "name": name},
		}

		result, err := securityQueryFields.Resolve(params.ResolveParams)

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
