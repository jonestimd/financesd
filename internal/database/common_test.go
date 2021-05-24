package database

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jonestimd/financesd/internal/database/table"
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

		result := runQuery(tx, companyType, query, name)

		assert.Equal(t, result, []*table.Company{{ID: 42, Name: name}})
	})
}

func Test_runQuery_panicsForQueryError(t *testing.T) {
	name := "the company"
	query := "select * from company where name = ?"
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		expectedErr := errors.New("query error")
		mock.ExpectQuery(query).WithArgs(name).WillReturnError(expectedErr)
		defer func() {
			if err := recover(); err != nil {
				assert.Same(t, expectedErr, err)
			} else {
				assert.Fail(t, "expected an error")
			}
		}()

		runQuery(tx, companyType, query, name)
	})
}

func Test_runUpdate_closesStatement(t *testing.T) {
	query := "update company set name = ? where id = ?"
	sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
		expectedArgs := []driver.Value{"new name", 42}
		rs := sqlmock.NewResult(-1, 1)
		mockDB.ExpectPrepare(query).WillBeClosed().ExpectExec().WithArgs(expectedArgs...).WillReturnResult(rs)

		rowCount := runUpdate(tx, query, "new name", 42)

		assert.Equal(t, int64(1), rowCount)
		assert.Nil(t, mockDB.ExpectationsWereMet())
	})
}

func Test_runUpdate_panicsForPrepareError(t *testing.T) {
	query := "update company set name = ? where id = ?"
	sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
		expectedErr := errors.New("query error")
		mockDB.ExpectPrepare(query).WillReturnError(expectedErr)
		defer func() {
			assert.Nil(t, mockDB.ExpectationsWereMet())
			if err := recover(); err != nil {
				assert.Same(t, expectedErr, err)
			} else {
				assert.Fail(t, "expected an error")
			}
		}()

		runUpdate(tx, query, "new name", 42)
	})
}

func Test_runUpdate_panicsForQueryError(t *testing.T) {
	query := "update company set name = ? where id = ?"
	sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
		expectedArgs := []driver.Value{"new name", 42}
		expectedErr := errors.New("query error")
		mockDB.ExpectPrepare(query).WillBeClosed().ExpectExec().WithArgs(expectedArgs...).WillReturnError(expectedErr)
		defer func() {
			assert.Nil(t, mockDB.ExpectationsWereMet())
			if err := recover(); err != nil {
				assert.Same(t, expectedErr, err)
			} else {
				assert.Fail(t, "expected an error")
			}
		}()

		runUpdate(tx, query, "new name", 42)
	})
}

func Test_runInsert_returnsID(t *testing.T) {
	id := int64(42)
	name := "the company"
	query := "insert into company (name) values(?)"
	sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
		prepare := mockDB.ExpectPrepare(query)
		prepare.ExpectExec().
			WithArgs(name).
			WillReturnResult(sqlmock.NewResult(id, 1))
		prepare.WillBeClosed()

		result := runInsert(tx, query, name)

		assert.Equal(t, id, result)
		assert.Nil(t, mockDB.ExpectationsWereMet())
	})
}

func Test_runInsert_panicsForPrepareError(t *testing.T) {
	expectedErr := errors.New("prepare failed")
	query := "insert into company (name) values(?)"
	sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
		mockDB.ExpectPrepare(query).WillReturnError(expectedErr)
		defer func() {
			assert.Nil(t, mockDB.ExpectationsWereMet())
			if err := recover(); err != nil {
				assert.Same(t, expectedErr, err)
			} else {
				assert.Fail(t, "expected an error")
			}
		}()

		runInsert(tx, query, "new company")
	})
}

func Test_runInsert_panicsForExecError(t *testing.T) {
	expectedErr := errors.New("query failed")
	name := "the company"
	query := "insert into company (name) values(?)"
	sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
		prepare := mockDB.ExpectPrepare(query)
		prepare.ExpectExec().WithArgs(name).WillReturnError(expectedErr)
		prepare.WillBeClosed()
		defer func() {
			assert.Nil(t, mockDB.ExpectationsWereMet())
			if err := recover(); err != nil {
				assert.Same(t, expectedErr, err)
			} else {
				assert.Fail(t, "expected an error")
			}
		}()

		runInsert(tx, query, name)
	})
}
