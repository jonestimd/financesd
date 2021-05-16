package database

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MonsantoCo/mocka/v2"
	"github.com/jonestimd/financesd/internal/sqltest"
	"github.com/stretchr/testify/assert"
)

func assertCompaniesError(t *testing.T, companies []*Company, err error, expectedErr error) {
	assert.Same(t, expectedErr, err)
	assert.Nil(t, companies)
}

func Test_GetAllCompanies(t *testing.T) {
	t.Run("returns companies", func(t *testing.T) {
		companies := []*Company{{ID: 1}}
		runQueryStub := mocka.Function(t, &runQuery, companies, nil)
		defer runQueryStub.Restore()
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result, err := GetAllCompanies(tx)

			assert.Equal(t, []interface{}{tx, companyType, "select * from company", []interface{}(nil)}, runQueryStub.GetFirstCall().Arguments())
			assert.Nil(t, err)
			assert.Equal(t, companies, result)
		})
	})
	t.Run("returns error", func(t *testing.T) {
		expectedErr := errors.New("database error")
		runQueryStub := mocka.Function(t, &runQuery, nil, expectedErr)
		defer runQueryStub.Restore()
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result, err := GetAllCompanies(tx)

			assertCompaniesError(t, result, err, expectedErr)
		})
	})
}

func Test_GetCompanyByID(t *testing.T) {
	id := int64(42)
	t.Run("returns companies", func(t *testing.T) {
		companies := []*Company{{ID: 1}}
		runQueryStub := mocka.Function(t, &runQuery, companies, nil)
		defer runQueryStub.Restore()
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result, err := GetCompanyByID(tx, id)

			assert.Equal(t, []interface{}{tx, companyType, "select * from company where id = ?", []interface{}{id}},
				runQueryStub.GetFirstCall().Arguments())
			assert.Nil(t, err)
			assert.Equal(t, companies, result)
		})
	})
	t.Run("returns error", func(t *testing.T) {
		expectedErr := errors.New("database error")
		runQueryStub := mocka.Function(t, &runQuery, nil, expectedErr)
		defer runQueryStub.Restore()
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result, err := GetCompanyByID(tx, id)

			assertCompaniesError(t, result, err, expectedErr)
		})
	})
}

func Test_GetCompanyByName(t *testing.T) {
	name := "the company"
	t.Run("returns companies", func(t *testing.T) {
		companies := []*Company{{ID: 1}}
		runQueryStub := mocka.Function(t, &runQuery, companies, nil)
		defer runQueryStub.Restore()
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result, err := GetCompanyByName(tx, name)

			assert.Equal(t, []interface{}{tx, companyType, "select * from company where name = ?", []interface{}{name}},
				runQueryStub.GetFirstCall().Arguments())
			assert.Nil(t, err)
			assert.Equal(t, companies, result)
		})
	})
	t.Run("returns error", func(t *testing.T) {
		expectedErr := errors.New("database error")
		runQueryStub := mocka.Function(t, &runQuery, nil, expectedErr)
		defer runQueryStub.Restore()
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result, err := GetCompanyByName(tx, name)

			assertCompaniesError(t, result, err, expectedErr)
		})
	})
}

func Test_GetCompaniesByIDs(t *testing.T) {
	ids := []int64{42, 96}
	t.Run("returns companies", func(t *testing.T) {
		companies := []*Company{{ID: 42}}
		runQueryStub := mocka.Function(t, &runQuery, companies, nil)
		defer runQueryStub.Restore()
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result, err := GetCompaniesByIDs(tx, ids)

			assert.Equal(t, []interface{}{tx, companyType, "select * from company where json_contains(?, cast(id as json))", []interface{}{"[42,96]"}},
				runQueryStub.GetFirstCall().Arguments())
			assert.Nil(t, err)
			assert.Equal(t, companies, result)
		})
	})
	t.Run("returns error", func(t *testing.T) {
		expectedErr := errors.New("database error")
		runQueryStub := mocka.Function(t, &runQuery, nil, expectedErr)
		defer runQueryStub.Restore()
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result, err := GetCompaniesByIDs(tx, ids)

			assert.Same(t, expectedErr, err)
			assert.Nil(t, result)
		})
	})
}

func Test_AddCompany(t *testing.T) {
	id := int64(42)
	t.Run("returns companies with ids", func(t *testing.T) {
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			runInsertStub := mocka.Function(t, &runInsert, id, nil)
			runInsertStub.OnSecondCall().Return(id+1, nil)
			defer runInsertStub.Restore()

			results, err := AddCompany(tx, "company1", "somebody")

			assert.Nil(t, err)
			assert.Equal(t, id, results.ID)
			assert.Equal(t, "company1", results.Name)
			assert.Equal(t, "somebody", results.ChangeUser)
			assert.Equal(t, 1, results.Version)
			assert.Nil(t, mock.ExpectationsWereMet())
		})
	})
	t.Run("returns error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			expectedErr := errors.New("database error")
			runInsertStub := mocka.Function(t, &runInsert, int64(0), expectedErr)
			defer runInsertStub.Restore()

			result, err := AddCompany(tx, "company1", "somebody")

			assert.Same(t, expectedErr, err)
			assert.Nil(t, result)
		})
	})
}

func Test_DeleteCompanies(t *testing.T) {
	ids := []int{42, 96}
	t.Run("returns update error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			expectedErr := errors.New("database error")
			runUpdateStub := mocka.Function(t, &runUpdate, int64(0), expectedErr)
			defer runUpdateStub.Restore()

			_, err := DeleteCompanies(tx, ids)

			assert.Same(t, expectedErr, err)
		})
	})
	t.Run("returns delete count", func(t *testing.T) {
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			count := int64(2)
			runUpdateStub := mocka.Function(t, &runUpdate, count, nil)
			defer runUpdateStub.Restore()

			result, err := DeleteCompanies(tx, ids)

			assert.Equal(t,
				[]interface{}{tx, "delete from company where json_contains(?, cast(id as json))", []interface{}{intsToJson(ids)}},
				runUpdateStub.GetCall(0).Arguments())
			assert.Equal(t, count, result)
			assert.Nil(t, err)
		})
	})
}

func Test_UpdateCompanies(t *testing.T) {
	tests := []struct {
		name        string
		count       int64
		updateError error
		message     string
	}{
		{"returns update error", 0, errors.New("database error"), "database error"},
		{"returns not found error", 0, nil, "company not found (42 @ 1)"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
				runUpdateStub := mocka.Function(t, &runUpdate, test.count, test.updateError)
				defer runUpdateStub.Restore()

				err := UpdateCompany(tx, 42, 1, "name", "somebody")

				assert.EqualError(t, err, test.message)
			})
		})
	}
	t.Run("returns companies", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			runUpdateStub := mocka.Function(t, &runUpdate, int64(1), nil)
			defer runUpdateStub.Restore()

			err := UpdateCompany(tx, 42, 1, "rename 42", "somebody")

			assert.Nil(t, err)
			assert.Equal(t, []interface{}{tx, updateCompanySQL, []interface{}{"rename 42", "somebody", int64(42), int64(1)}},
				runUpdateStub.GetCall(0).Arguments())
		})
	})
}
