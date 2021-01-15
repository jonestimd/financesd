package graphql

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
)

func Test_groupQueryFields_Resolve(t *testing.T) {
	testQuery(t, func(mock sqlmock.Sqlmock, orm *gorm.DB) {
		params := newResolveParams(orm, groupQuery, nil, newField("", "id"), newField("", "name"))
		rows := sqlmock.NewRows([]string{"json"}).AddRow(`{"id":1,"name":"Group 1"}`)
		mock.ExpectQuery(`^select json_object\("id", g\.id, "name", g\.name\) from tx_group g$`).WithArgs().WillReturnRows(rows)
		expectedRows := []interface{}{
			map[string]interface{}{"id": 1.0, "name": "Group 1"},
		}

		result, err := groupQueryFields.Resolve(params)

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
