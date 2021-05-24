package database

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MonsantoCo/mocka/v2"
	"github.com/jonestimd/financesd/internal/database/table"
	"github.com/jonestimd/financesd/internal/sqltest"
	"github.com/stretchr/testify/assert"
)

func Test_GetAllCompanies(t *testing.T) {
	companies := []*table.Company{{ID: 1}}
	runQueryStub := mocka.Function(t, &runQuery, companies)
	defer runQueryStub.Restore()
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		result := GetAllCompanies(tx)

		assert.Equal(t, []interface{}{tx, companyType, "select * from company", []interface{}(nil)}, runQueryStub.GetFirstCall().Arguments())
		assert.Equal(t, companies, result)
	})
}

func Test_GetCompanyByID(t *testing.T) {
	id := int64(42)
	companies := []*table.Company{{ID: 1}}
	runQueryStub := mocka.Function(t, &runQuery, companies)
	defer runQueryStub.Restore()
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		result := GetCompanyByID(tx, id)

		assert.Equal(t, []interface{}{tx, companyType, "select * from company where id = ?", []interface{}{id}},
			runQueryStub.GetFirstCall().Arguments())
		assert.Equal(t, companies, result)
	})
}

func Test_GetCompanyByName(t *testing.T) {
	name := "the company"
	companies := []*table.Company{{ID: 1}}
	runQueryStub := mocka.Function(t, &runQuery, companies)
	defer runQueryStub.Restore()
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		result := GetCompanyByName(tx, name)

		assert.Equal(t, []interface{}{tx, companyType, "select * from company where name = ?", []interface{}{name}},
			runQueryStub.GetFirstCall().Arguments())
		assert.Equal(t, companies, result)
	})
}

func Test_GetCompaniesByIDs(t *testing.T) {
	ids := []int64{42, 96}
	companies := []*table.Company{{ID: 42}}
	runQueryStub := mocka.Function(t, &runQuery, companies)
	defer runQueryStub.Restore()
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		result := GetCompaniesByIDs(tx, ids)

		assert.Equal(t, []interface{}{tx, companyType, "select * from company where json_contains(?, cast(id as json))", []interface{}{"[42,96]"}},
			runQueryStub.GetFirstCall().Arguments())
		assert.Equal(t, companies, result)
	})
}

func Test_AddCompany(t *testing.T) {
	id := int64(42)
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		runInsertStub := mocka.Function(t, &runInsert, id)
		runInsertStub.OnSecondCall().Return(id + 1)
		defer runInsertStub.Restore()

		results := AddCompany(tx, "company1", "somebody")

		assert.Equal(t, id, results.ID)
		assert.Equal(t, "company1", results.Name)
		assert.Equal(t, "somebody", results.ChangeUser)
		assert.Equal(t, 1, results.Version)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func Test_DeleteCompanies(t *testing.T) {
	ids := []int{42, 96}
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		count := int64(2)
		runUpdateStub := mocka.Function(t, &runUpdate, count)
		defer runUpdateStub.Restore()

		result := DeleteCompanies(tx, ids)

		assert.Equal(t,
			[]interface{}{tx, "delete from company where json_contains(?, cast(id as json))", []interface{}{intsToJson(ids)}},
			runUpdateStub.GetCall(0).Arguments())
		assert.Equal(t, count, result)
	})
}

func Test_UpdateCompanies(t *testing.T) {
	t.Run("panics if not found", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			defer func() {
				if err := recover(); err != nil {
					assert.Equal(t, "company not found (42 @ 1)", err.(error).Error())
				} else {
					assert.Fail(t, "expected an error")
				}
			}()
			runUpdateStub := mocka.Function(t, &runUpdate, int64(0))
			defer runUpdateStub.Restore()

			UpdateCompany(tx, 42, 1, "name", "somebody")
		})
	})
	t.Run("returns companies", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			runUpdateStub := mocka.Function(t, &runUpdate, int64(1))
			defer runUpdateStub.Restore()

			UpdateCompany(tx, 42, 1, "rename 42", "somebody")

			assert.Equal(t, []interface{}{tx, updateCompanySQL, []interface{}{"rename 42", "somebody", int64(42), int64(1)}},
				runUpdateStub.GetCall(0).Arguments())
		})
	})
}
