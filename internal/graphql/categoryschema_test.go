package graphql

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
)

func Test_categoryQueryFields_Resolve(t *testing.T) {
	testQuery(t, func(mock sqlmock.Sqlmock, orm *gorm.DB) {
		params := newResolveParams(orm, categoryQuery, nil, newField("", "id"), newField("", "code"))
		rows := sqlmock.NewRows([]string{"json"}).AddRow(`{"id":1,"code":"Category 1"}`)
		mock.ExpectQuery(`^select json_object\("id", c\.id, "code", c\.code\) from transaction_category c$`).WithArgs().WillReturnRows(rows)
		expectedRows := []interface{}{
			map[string]interface{}{"id": 1.0, "code": "Category 1"},
		}

		result, err := categoryQueryFields.Resolve(params)

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
