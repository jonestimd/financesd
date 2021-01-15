package graphql

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
)

func testQuery(t *testing.T, test func(mockDb sqlmock.Sqlmock, orm *gorm.DB)) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	orm, err := gorm.Open("postgres", db) // TODO MySql not supported?
	if err != nil {
		panic(err.Error())
	}
	defer orm.Close()
	orm.SingularTable(true)

	test(mock, orm)
}

func mockRows(columns ...string) *sqlmock.Rows {
	return sqlmock.NewRows(columns)
}
