package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MonsantoCo/mocka/v2"
	"github.com/jonestimd/financesd/internal/sqltest"
	"github.com/stretchr/testify/assert"
)

func Test_company_ptrTo(t *testing.T) {
	company := &Company{}
	tests := []struct {
		column string
		ptr    interface{}
	}{
		{column: "id", ptr: &company.ID},
		{column: "name", ptr: &company.Name},
		{column: "version", ptr: &company.Version},
		{column: "change_user", ptr: &company.ChangeUser},
		{column: "change_date", ptr: &company.ChangeDate},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s returns pointer to field", test.column), func(t *testing.T) {
			field := company.ptrTo(test.column)

			assert.Same(t, test.ptr, field)
		})
	}
}

func Test_Companies(t *testing.T) {
	companyID1 := int64(42)
	companyID2 := int64(69)
	companies := []*Company{{ID: companyID1}, {ID: companyID2}}
	c := &Companies{Companies: companies}
	t.Run("CompanyIDs returns IDs", func(t *testing.T) {
		companyIDs := c.CompanyIDs()

		assert.Equal(t, []int64{companyID1, companyID2}, companyIDs)
	})
	t.Run("Result returns companies", func(t *testing.T) {
		assert.Equal(t, c.Result(), companies)
	})
	t.Run("Error returns err", func(t *testing.T) {
		err := errors.New("test error")

		assert.Nil(t, c.Error())
		assert.Same(t, NewAccounts(err).Error(), err)
	})
}

func assertCompaniesError(t *testing.T, result *Companies, err error) {
	assert.Same(t, err, result.Error())
	assert.Nil(t, result.Companies)
}

func Test_GetAllCompanies(t *testing.T) {
	t.Run("returns companies", func(t *testing.T) {
		companies := []*Company{{ID: 1}}
		runQueryStub := mocka.Function(t, &runQuery, companies, nil)
		defer runQueryStub.Restore()
		sqltest.TestQuery(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result := GetAllCompanies(tx)

			assert.Equal(t, []interface{}{tx, companyType, "select * from company", []interface{}(nil)}, runQueryStub.GetFirstCall().Arguments())
			assert.Nil(t, result.Error())
			assert.Equal(t, companies, result.Companies)
		})
	})
	t.Run("returns error", func(t *testing.T) {
		err := errors.New("database error")
		runQueryStub := mocka.Function(t, &runQuery, nil, err)
		defer runQueryStub.Restore()
		sqltest.TestQuery(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result := GetAllCompanies(tx)

			assertCompaniesError(t, result, err)
		})
	})
}

func Test_GetCompanyByID(t *testing.T) {
	id := int64(42)
	t.Run("returns companies", func(t *testing.T) {
		companies := []*Company{{ID: 1}}
		runQueryStub := mocka.Function(t, &runQuery, companies, nil)
		defer runQueryStub.Restore()
		sqltest.TestQuery(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result := GetCompanyByID(tx, id)

			assert.Equal(t, []interface{}{tx, companyType, "select * from company where id = ?", []interface{}{id}},
				runQueryStub.GetFirstCall().Arguments())
			assert.Nil(t, result.Error())
			assert.Equal(t, companies, result.Companies)
		})
	})
	t.Run("returns error", func(t *testing.T) {
		err := errors.New("database error")
		runQueryStub := mocka.Function(t, &runQuery, nil, err)
		defer runQueryStub.Restore()
		sqltest.TestQuery(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result := GetCompanyByID(tx, id)

			assertCompaniesError(t, result, err)
		})
	})
}

func Test_GetCompanyByName(t *testing.T) {
	name := "the company"
	t.Run("returns companies", func(t *testing.T) {
		companies := []*Company{{ID: 1}}
		runQueryStub := mocka.Function(t, &runQuery, companies, nil)
		defer runQueryStub.Restore()
		sqltest.TestQuery(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result := GetCompanyByName(tx, name)

			assert.Equal(t, []interface{}{tx, companyType, "select * from company where name = ?", []interface{}{name}},
				runQueryStub.GetFirstCall().Arguments())
			assert.Nil(t, result.Error())
			assert.Equal(t, companies, result.Companies)
		})
	})
	t.Run("returns error", func(t *testing.T) {
		err := errors.New("database error")
		runQueryStub := mocka.Function(t, &runQuery, nil, err)
		defer runQueryStub.Restore()
		sqltest.TestQuery(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result := GetCompanyByName(tx, name)

			assertCompaniesError(t, result, err)
		})
	})
}

func Test_GetCompaniesByIDs(t *testing.T) {
	ids := []int64{42, 96}
	t.Run("returns companies", func(t *testing.T) {
		companies := []*Company{{ID: 42}}
		runQueryStub := mocka.Function(t, &runQuery, companies, nil)
		defer runQueryStub.Restore()
		sqltest.TestQuery(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result, err := GetCompaniesByIDs(tx, ids)

			jsonIDs, _ := json.Marshal(ids)
			assert.Equal(t, []interface{}{tx, companyType, "select * from company where json_contains(?, cast(id as json))", []interface{}{jsonIDs}},
				runQueryStub.GetFirstCall().Arguments())
			assert.Nil(t, err)
			assert.Equal(t, map[int64]*Company{companies[0].ID: companies[0]}, result)
		})
	})
	t.Run("returns error", func(t *testing.T) {
		expectedErr := errors.New("database error")
		runQueryStub := mocka.Function(t, &runQuery, nil, expectedErr)
		defer runQueryStub.Restore()
		sqltest.TestQuery(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			result, err := GetCompaniesByIDs(tx, ids)

			assert.Same(t, expectedErr, err)
			assert.Nil(t, result)
		})
	})
}
