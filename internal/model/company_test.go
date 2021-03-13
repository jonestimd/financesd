package model

import (
	"database/sql"
	"database/sql/driver"
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

func Test_AddCompanies(t *testing.T) {
	id := int64(42)
	t.Run("returns companies with ids", func(t *testing.T) {
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			runInsertStub := mocka.Function(t, &runInsert, id, nil)
			runInsertStub.OnSecondCall().Return(id+1, nil)
			defer runInsertStub.Restore()
			validateNameStub := mocka.Function(t, &validateName, nil)
			defer validateNameStub.Restore()

			results, err := AddCompanies(tx, []string{"company1", "company2"}, "somebody")

			assert.Nil(t, err)
			assert.Equal(t, 2, validateNameStub.CallCount())
			assert.Equal(t, []interface{}{"company1"}, validateNameStub.GetCall(0).Arguments())
			assert.Equal(t, []interface{}{"company2"}, validateNameStub.GetCall(1).Arguments())
			assert.Equal(t, id, results[0].ID)
			assert.Equal(t, "company1", results[0].Name)
			assert.Equal(t, "somebody", results[0].ChangeUser)
			assert.Equal(t, 1, results[0].Version)
			assert.NotNil(t, results[0].source)
			assert.Equal(t, id+1, results[1].ID)
			assert.Equal(t, "company2", results[1].Name)
			assert.Equal(t, "somebody", results[1].ChangeUser)
			assert.Equal(t, 1, results[1].Version)
			assert.NotNil(t, results[1].source)
			assert.Nil(t, mock.ExpectationsWereMet())
		})
	})
	t.Run("returns error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			expectedErr := errors.New("database error")
			runInsertStub := mocka.Function(t, &runInsert, int64(42), nil)
			runInsertStub.OnSecondCall().Return(int64(0), expectedErr)
			defer runInsertStub.Restore()

			result, err := AddCompanies(tx, []string{"company1", "company2"}, "somebody")

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
			runUpdateStub := mocka.Function(t, &runUpdate, nil, expectedErr)
			defer runUpdateStub.Restore()

			_, err := DeleteCompanies(tx, ids, "somebody")

			assert.Same(t, expectedErr, err)
		})
	})
	t.Run("returns delete count", func(t *testing.T) {
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			count := int64(2)
			runUpdateStub := mocka.Function(t, &runUpdate, sqlmock.NewResult(0, count), nil)
			defer runUpdateStub.Restore()

			result, err := DeleteCompanies(tx, ids, "somebody")

			assert.Equal(t,
				[]interface{}{tx, "delete from company where json_contains(?, cast(id as json))", []interface{}{intsToJson(ids)}},
				runUpdateStub.GetCall(0).Arguments())
			assert.Equal(t, count, result)
			assert.Nil(t, err)
		})
	})
}

func Test_UpdateCompanies(t *testing.T) {
	updates := []interface{}{
		map[string]interface{}{"id": 42, "name": "rename 42"},
		map[string]interface{}{"id": 96, "name": "rename 96"},
	}
	tests := []struct {
		name        string
		result      driver.Result
		updateError error
		getError    error
		message     string
	}{
		{"returns update error", nil, errors.New("database error"), nil, "database error"},
		{"returns count error", sqlmock.NewErrorResult(errors.New("database error")), nil, nil, "database error"},
		{"returns not found error", sqlmock.NewResult(0, 0), nil, nil, "company not found: 42"},
		{"returns getCompanies error", sqlmock.NewResult(0, 1), nil, errors.New("database error"), "database error"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
				runUpdateStub := mocka.Function(t, &runUpdate, test.result, test.updateError)
				defer runUpdateStub.Restore()
				getCompaniesStub := mocka.Function(t, &getCompaniesByIDs, nil, test.getError)
				defer getCompaniesStub.Restore()

				_, err := UpdateCompanies(tx, updates, "somebody")

				assert.EqualError(t, err, test.message)
			})
		})
	}
	t.Run("returns companies", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			companies := []*Company{{ID: 42, Name: "new name"}}
			runUpdateStub := mocka.Function(t, &runUpdate, sqlmock.NewResult(0, 1), nil)
			defer runUpdateStub.Restore()
			getCompaniesStub := mocka.Function(t, &getCompaniesByIDs, companies, nil)
			defer getCompaniesStub.Restore()

			result, err := UpdateCompanies(tx, updates, "somebody")

			assert.Equal(t, companies, result)
			assert.Nil(t, err)
			assert.Equal(t, 2, runUpdateStub.CallCount())
			assert.Equal(t, []interface{}{tx, "update company set name = ?, change_date = current_timestamp, change_user = ? where id = ?",
				[]interface{}{"rename 42", "somebody", int64(42)}}, runUpdateStub.GetCall(0).Arguments())
			assert.Equal(t, []interface{}{tx, "update company set name = ?, change_date = current_timestamp, change_user = ? where id = ?",
				[]interface{}{"rename 96", "somebody", int64(96)}}, runUpdateStub.GetCall(1).Arguments())
			assert.Equal(t, 1, getCompaniesStub.CallCount())
			assert.Equal(t, []interface{}{tx, []int64{42, 96}}, getCompaniesStub.GetCall(0).Arguments())
		})
	})
}
