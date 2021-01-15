package graphql

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
)

func Test_accountQueryFields_Resolve(t *testing.T) {
	testQuery(t, func(mock sqlmock.Sqlmock, orm *gorm.DB) {
		params := newResolveParams(orm, accountQuery, nil, newField("", "id"), newField("", "name"))
		rows := sqlmock.NewRows([]string{"json"}).AddRow(`{"id":1,"name":"Account 1"}`)
		mock.ExpectQuery(`^select json_object\("id", a\.id, "name", a\.name\) from account a$`).WithArgs().WillReturnRows(rows)
		expectedRows := []interface{}{
			map[string]interface{}{"id": 1.0, "name": "Account 1"},
		}

		result, err := accountQueryFields.Resolve(params)

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
