package model

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jonestimd/financesd/internal/sqltest"
	"github.com/stretchr/testify/assert"
)

func Test_intsToJson(t *testing.T) {
	values := []int{1, 3, 5, 2, 4, 6}

	result := intsToJson(values)

	assert.Equal(t, result, "[1,3,5,2,4,6]")
}

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

func Test_runUpdate_closesStatement(t *testing.T) {
	query := "update company set name = ? where id = ?"
	sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
		expectedArgs := []driver.Value{"new name", 42}
		rs := sqlmock.NewResult(-1, 1)
		mockDB.ExpectPrepare(query).WillBeClosed().ExpectExec().WithArgs(expectedArgs...).WillReturnResult(rs)

		rowCount, err := runUpdate(tx, query, "new name", 42)

		assert.Nil(t, err)
		assert.Equal(t, int64(1), rowCount)
		assert.Nil(t, mockDB.ExpectationsWereMet())
	})
}

func Test_runUpdate_returnsPrepareError(t *testing.T) {
	query := "update company set name = ? where id = ?"
	sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
		expectedErr := errors.New("query error")
		mockDB.ExpectPrepare(query).WillReturnError(expectedErr)

		count, err := runUpdate(tx, query, "new name", 42)

		assert.Same(t, expectedErr, err)
		assert.Equal(t, int64(0), count)
		assert.Nil(t, mockDB.ExpectationsWereMet())
	})
}

func Test_runUpdate_returnsQueryError(t *testing.T) {
	query := "update company set name = ? where id = ?"
	sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
		expectedArgs := []driver.Value{"new name", 42}
		expectedErr := errors.New("query error")
		mockDB.ExpectPrepare(query).WillBeClosed().ExpectExec().WithArgs(expectedArgs...).WillReturnError(expectedErr)

		count, err := runUpdate(tx, query, "new name", 42)

		assert.Same(t, expectedErr, err)
		assert.Equal(t, int64(0), count)
		assert.Nil(t, mockDB.ExpectationsWereMet())
	})
}

func Test_runInsert_returnsID(t *testing.T) {
	id := int64(42)
	name := "the company"
	query := "insert into company (name) values(?)"
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		prepare := mock.ExpectPrepare(query)
		prepare.ExpectExec().
			WithArgs(name).
			WillReturnResult(sqlmock.NewResult(id, 1))
		prepare.WillBeClosed()

		result, err := runInsert(tx, query, name)

		assert.Nil(t, err)
		assert.Equal(t, id, result)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func Test_runInsert_returnsPrepareError(t *testing.T) {
	expectedErr := errors.New("prepare failed")
	query := "insert into company (name) values(?)"
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		mock.ExpectPrepare(query).WillReturnError(expectedErr)

		result, err := runInsert(tx, query, "new company")

		assert.Equal(t, int64(0), result)
		assert.Same(t, expectedErr, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func Test_runInsert_returnsExecError(t *testing.T) {
	expectedErr := errors.New("query failed")
	name := "the company"
	query := "insert into company (name) values(?)"
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		prepare := mock.ExpectPrepare(query)
		prepare.ExpectExec().
			WithArgs(name).
			WillReturnError(expectedErr)
		prepare.WillBeClosed()

		result, err := runInsert(tx, query, name)

		assert.Equal(t, int64(0), result)
		assert.Same(t, expectedErr, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}
