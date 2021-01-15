package graphql

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	"github.com/jonestimd/financesd/internal/model"
)

func Test_securitySchema_includesAssetFields(t *testing.T) {
	assetFields := map[string]string{
		"id":         "ID",
		"name":       "Name",
		"scale":      "Scale",
		"symbol":     "Symbol",
		"version":    "Version",
		"changeUser": "ChangeUser",
		"changeDate": "ChangeDate",
	}
	symbol := "A1"
	now := time.Now()
	asset := model.Asset{
		ID:      1,
		Name:    "Asset 1",
		Scale:   6,
		Symbol:  &symbol,
		Version: 9,
		Audited: model.Audited{
			ChangeUser: "me",
			ChangeDate: &now,
		},
	}
	security := model.Security{
		Type:  "Stock",
		Asset: asset,
	}

	securityFields := securitySchema.Fields()

	if len(securityFields) != len(assetFields)+1 {
		t.Errorf("Expected %v fields, got %v", len(assetFields)+1, len(securityFields))
	}
	resolveParams := graphql.ResolveParams{Source: &security}
	assetAccessor := reflect.Indirect(reflect.ValueOf(&asset))
	for name, structName := range assetFields {
		field := securityFields[name]
		if field == nil {
			t.Errorf("Missing field: %s", name)
		}
		value, err := field.Resolve(resolveParams)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		expectedValue := assetAccessor.FieldByName(structName).Interface()
		if value != expectedValue {
			t.Errorf("Expected %v for %s, got %v", expectedValue, name, value)
		}
	}
}

func Test_securityQueryFields_Resolve_all(t *testing.T) {
	testQuery(t, func(mock sqlmock.Sqlmock, orm *gorm.DB) {
		var assetID int64 = 1
		params := newResolveParams(orm, securityQuery, nil, newField("", "type"), newField("", "id"), newField("", "name"))
		mock.ExpectQuery(`^SELECT \* FROM "security"$`).WithArgs().WillReturnRows(
			mockRows("asset_id", "type").AddRow(assetID, "Stock"))
		mock.ExpectQuery(`^SELECT \* FROM "asset" WHERE \("id" IN \(\$1\)\)$`).WithArgs(assetID).WillReturnRows(
			mockRows("id", "name").AddRow(assetID, "Security Name"))
		expectedRows := []*model.Security{
			{AssetID: assetID, Asset: model.Asset{ID: assetID, Name: "Security Name"}, Type: "Stock"},
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

func Test_securityQueryFields_Resolve_byID(t *testing.T) {
	testQuery(t, func(mock sqlmock.Sqlmock, orm *gorm.DB) {
		var assetID int64 = 1
		args := map[string]interface{}{"id": fmt.Sprint(assetID)}
		params := newResolveParams(orm, securityQuery, args, newField("", "type"), newField("", "id"), newField("", "name"))
		mock.ExpectQuery(`^SELECT \* FROM "security" WHERE \(asset_id = \$1\)$`).WithArgs(assetID).WillReturnRows(
			mockRows("asset_id", "type").AddRow(assetID, "Stock"))
		mock.ExpectQuery(`^SELECT \* FROM "asset" WHERE \("id" IN \(\$1\)\)$`).WithArgs(assetID).WillReturnRows(
			mockRows("id", "name").AddRow(assetID, "Security Name"))
		expectedRows := []*model.Security{
			{AssetID: assetID, Asset: model.Asset{ID: assetID, Name: "Security Name"}, Type: "Stock"},
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
	testQuery(t, func(mock sqlmock.Sqlmock, orm *gorm.DB) {
		args := map[string]interface{}{"id": "one"}
		params := newResolveParams(orm, securityQuery, args, newField("", "type"), newField("", "id"), newField("", "name"))

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
	testQuery(t, func(mock sqlmock.Sqlmock, orm *gorm.DB) {
		var assetID int64 = 1
		symbol := "A1"
		args := map[string]interface{}{"symbol": symbol}
		params := newResolveParams(orm, securityQuery, args, newField("", "type"), newField("", "id"), newField("", "name"))
		mock.ExpectQuery(`^SELECT "security"\.\* FROM "security" join asset a on a\.id = security\.asset_id WHERE \(a\.symbol = \$1\)$`).
			WithArgs(symbol).
			WillReturnRows(mockRows("asset_id", "type").AddRow(assetID, "Stock"))
		mock.ExpectQuery(`^SELECT \* FROM "asset" WHERE \("id" IN \(\$1\)\)$`).WithArgs(assetID).WillReturnRows(
			mockRows("id", "name").AddRow(assetID, "Security Name"))
		expectedRows := []*model.Security{
			{AssetID: assetID, Asset: model.Asset{ID: assetID, Name: "Security Name"}, Type: "Stock"},
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
