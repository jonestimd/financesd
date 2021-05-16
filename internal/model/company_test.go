package model

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MonsantoCo/mocka/v2"
	"github.com/jonestimd/financesd/internal/database"
	"github.com/jonestimd/financesd/internal/sqltest"
	"github.com/stretchr/testify/assert"
)

func Test_Company_GetAccounts(t *testing.T) {
	companyID := int64(42)
	companyID2 := int64(96)
	account1 := &Account{Account: &database.Account{CompanyID: &companyID}}
	account2 := &Account{Account: &database.Account{CompanyID: &companyID2}}
	company := &database.Company{ID: companyID}
	loadedSource := &companySource{accounts: []*Account{account1, account2}}
	expectedErr := errors.New("test error")
	tests := []struct {
		name     string
		company  *Company
		accounts []*Account
		err      error
	}{
		{"returns empty array for no accounts", &Company{Company: company, source: &companySource{accounts: []*Account{}}}, []*Account{}, nil},
		{"returns existing accounts", &Company{Company: company, source: loadedSource}, []*Account{account1}, nil},
		{"returns existing source error", &Company{Company: company, source: &companySource{err: expectedErr}}, nil, expectedErr},
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
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			companies := []*database.Company{{ID: 1}}
			getAllCompaniesStub := mocka.Function(t, &getAllCompanies, companies, nil)
			defer getAllCompaniesStub.Restore()

			result, err := GetAllCompanies(tx)

			assert.Equal(t, []interface{}{tx}, getAllCompaniesStub.GetFirstCall().Arguments())
			assert.Nil(t, err)
			assert.Equal(t, companies[0], result[0].Company)
			assert.NotNil(t, result[0].source)
		})
	})
	t.Run("returns error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			expectedErr := errors.New("database error")
			getAllCompaniesStub := mocka.Function(t, &getAllCompanies, nil, expectedErr)
			defer getAllCompaniesStub.Restore()

			result, err := GetAllCompanies(tx)

			assertCompaniesError(t, result, err, expectedErr)
		})
	})
}

func Test_GetCompanyByID(t *testing.T) {
	id := int64(42)
	t.Run("returns companies", func(t *testing.T) {
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			companies := []*database.Company{{ID: 1}}
			getCompanyByIDStub := mocka.Function(t, &getCompanyByID, companies, nil)
			defer getCompanyByIDStub.Restore()

			result, err := GetCompanyByID(tx, id)

			assert.Equal(t, []interface{}{tx, id}, getCompanyByIDStub.GetFirstCall().Arguments())
			assert.Nil(t, err)
			assert.Equal(t, companies[0], result[0].Company)
			assert.NotNil(t, result[0].source)
		})
	})
	t.Run("returns error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			expectedErr := errors.New("database error")
			getCompanyByIDStub := mocka.Function(t, &getCompanyByID, nil, expectedErr)
			defer getCompanyByIDStub.Restore()

			result, err := GetCompanyByID(tx, id)

			assertCompaniesError(t, result, err, expectedErr)
		})
	})
}

func Test_GetCompanyByName(t *testing.T) {
	name := "the company"
	t.Run("returns companies", func(t *testing.T) {
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			companies := []*database.Company{{ID: 1}}
			getCompanyByNameStub := mocka.Function(t, &getCompanyByName, companies, nil)
			defer getCompanyByNameStub.Restore()

			result, err := GetCompanyByName(tx, name)

			assert.Equal(t, []interface{}{tx, name}, getCompanyByNameStub.GetFirstCall().Arguments())
			assert.Nil(t, err)
			assert.Equal(t, companies[0], result[0].Company)
			assert.NotNil(t, result[0].source)
		})
	})
	t.Run("returns error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			expectedErr := errors.New("database error")
			getCompanyByNameStub := mocka.Function(t, &getCompanyByName, nil, expectedErr)
			defer getCompanyByNameStub.Restore()

			result, err := GetCompanyByName(tx, name)

			assertCompaniesError(t, result, err, expectedErr)
		})
	})
}

func Test_GetCompaniesByIDs(t *testing.T) {
	ids := []int64{42, 96}
	t.Run("returns companies", func(t *testing.T) {
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			companies := []*database.Company{{ID: 42}}
			getCompaniesByIDsStub := mocka.Function(t, &getCompaniesByIDs, companies, nil)
			defer getCompaniesByIDsStub.Restore()

			result, err := GetCompaniesByIDs(tx, ids)

			assert.Equal(t, []interface{}{tx, ids}, getCompaniesByIDsStub.GetFirstCall().Arguments())
			assert.Nil(t, err)
			assert.Equal(t, companies[0], result[0].Company)
			assert.NotNil(t, result[0].source)
		})
	})
	t.Run("returns error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			expectedErr := errors.New("database error")
			getCompaniesByIDsStub := mocka.Function(t, &getCompaniesByIDs, nil, expectedErr)
			defer getCompaniesByIDsStub.Restore()

			result, err := GetCompaniesByIDs(tx, ids)

			assert.Same(t, expectedErr, err)
			assert.Nil(t, result)
		})
	})
}

func Test_AddCompanies(t *testing.T) {
	id := int64(42)
	t.Run("returns companies with ids", func(t *testing.T) {
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			company1 := &database.Company{ID: id}
			company2 := &database.Company{ID: id + 1}
			addCompanyStub := mocka.Function(t, &addCompany, company1, nil)
			addCompanyStub.OnCall(1).Return(company2, nil)
			defer addCompanyStub.Restore()
			validateNameStub := mocka.Function(t, &validateName, nil)
			defer validateNameStub.Restore()

			results, err := AddCompanies(tx, []string{"company1", "company2"}, "somebody")

			assert.Nil(t, err)
			assert.Equal(t, 2, validateNameStub.CallCount())
			assert.Equal(t, []interface{}{"company1"}, validateNameStub.GetCall(0).Arguments())
			assert.Equal(t, []interface{}{"company2"}, validateNameStub.GetCall(1).Arguments())
			assert.Equal(t, company1, results[0].Company)
			assert.NotNil(t, results[0].source)
			assert.Equal(t, company2, results[1].Company)
			assert.NotNil(t, results[1].source)
		})
	})
	t.Run("returns error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
			expectedErr := errors.New("database error")
			addCompanyStub := mocka.Function(t, &addCompany, nil, expectedErr)
			defer addCompanyStub.Restore()

			result, err := AddCompanies(tx, []string{"company1"}, "somebody")

			assert.Same(t, expectedErr, err)
			assert.Nil(t, result)
		})
	})
}

func Test_UpdateCompanies(t *testing.T) {
	updates := []interface{}{
		map[string]interface{}{"id": 42, "name": "rename 42", "version": 1},
		map[string]interface{}{"id": 96, "name": "rename 96", "version": 1},
	}
	tests := []struct {
		name        string
		updateError error
		getError    error
		message     string
	}{
		{"returns update error", errors.New("database error"), nil, "database error"},
		{"returns getCompanies error", nil, errors.New("database error"), "database error"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
				updateCompanyStub := mocka.Function(t, &updateCompany, test.updateError)
				defer updateCompanyStub.Restore()
				getCompaniesStub := mocka.Function(t, &GetCompaniesByIDs, nil, test.getError)
				defer getCompaniesStub.Restore()

				_, err := UpdateCompanies(tx, updates, "somebody")

				assert.EqualError(t, err, test.message)
			})
		})
	}
	t.Run("returns companies", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			companies := []*Company{{Company: &database.Company{ID: 42, Name: "new name"}}}
			updateCompanyStub := mocka.Function(t, &updateCompany, nil)
			defer updateCompanyStub.Restore()
			getCompaniesStub := mocka.Function(t, &GetCompaniesByIDs, companies, nil)
			defer getCompaniesStub.Restore()

			result, err := UpdateCompanies(tx, updates, "somebody")

			assert.Equal(t, companies, result)
			assert.Nil(t, err)
			assert.Equal(t, 2, updateCompanyStub.CallCount())
			assert.Equal(t, []interface{}{tx, int64(42), int64(1), "rename 42", "somebody"}, updateCompanyStub.GetCall(0).Arguments())
			assert.Equal(t, []interface{}{tx, int64(96), int64(1), "rename 96", "somebody"}, updateCompanyStub.GetCall(1).Arguments())
			assert.Equal(t, 1, getCompaniesStub.CallCount())
			assert.Equal(t, []interface{}{tx, []int64{42, 96}}, getCompaniesStub.GetCall(0).Arguments())
		})
	})
}
