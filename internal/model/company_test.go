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

func Test_Company_GetAccounts(t *testing.T) {
	companyID := int64(42)
	companyID2 := int64(96)
	account1 := &Account{CompanyID: &companyID}
	account2 := &Account{CompanyID: &companyID2}
	loadedSource := &companySource{accounts: []*Account{account1, account2}}
	expectedErr := errors.New("test error")
	tests := []struct {
		name     string
		company  *Company
		accounts []*Account
		err      error
	}{
		{"returns empty array for no accounts", &Company{ID: companyID, source: &companySource{accounts: []*Account{}}}, []*Account{}, nil},
		{"returns existing accounts", &Company{ID: companyID, source: loadedSource}, []*Account{account1}, nil},
		{"returns existing source error", &Company{ID: companyID, source: &companySource{err: expectedErr}}, nil, expectedErr},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := test.company.GetAccounts(nil)

			assert.Equal(t, test.err, err)
			assert.Equal(t, test.accounts, result)
		})
	}
}

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
			result, err := getCompaniesByIDs(tx, ids)

			jsonIDs, _ := json.Marshal(ids)
			assert.Equal(t, []interface{}{tx, companyType, "select * from company where json_contains(?, cast(id as json))", []interface{}{jsonIDs}},
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
			result, err := getCompaniesByIDs(tx, ids)

			assert.Same(t, expectedErr, err)
			assert.Nil(t, result)
		})
	})
}