package model

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jonestimd/financesd/internal/sqltest"
	"github.com/stretchr/testify/assert"
)

func Test_runQuery_populatesModel(t *testing.T) {
	name := "the company"
	query := "select * from company where name = ?"
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		rows := sqltest.MockRows("id", "name").AddRow(42, name)
		mock.ExpectQuery(query).WithArgs(name).WillReturnRows(rows)

		result, err := runQuery(tx, companyType, query, name)

		assert.Nil(t, err)
		assert.Equal(t, result, []*Company{{ID: 42, Name: name}})
	})
}

func Test_runQuery_returnsQueryError(t *testing.T) {
	name := "the company"
	query := "select * from company where name = ?"
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		expectedErr := errors.New("query error")
		mock.ExpectQuery(query).WithArgs(name).WillReturnError(expectedErr)

		result, err := runQuery(tx, companyType, query, name)

		assert.Same(t, expectedErr, err)
		assert.Nil(t, result)
	})
}
