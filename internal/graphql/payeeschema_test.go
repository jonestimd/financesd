package graphql

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
)

func Test_payeeQueryFields_Resolve(t *testing.T) {
	testQuery(t, func(mock sqlmock.Sqlmock, orm *gorm.DB) {
		params := newResolveParams(orm, payeeQuery, nil, newField("", "id"), newField("", "name"))
		rows := sqlmock.NewRows([]string{"json"}).AddRow(`{"id":1,"name":"Payee 1"}`)
		mock.ExpectQuery(`^select json_object\("id", p\.id, "name", p\.name\) from payee p$`).WithArgs().WillReturnRows(rows)
		expectedRows := []interface{}{
			map[string]interface{}{"id": 1.0, "name": "Payee 1"},
		}

		result, err := payeeQueryFields.Resolve(params)

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
