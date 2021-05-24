package domain

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MonsantoCo/mocka/v2"
	"github.com/jonestimd/financesd/internal/database/table"
	"github.com/jonestimd/financesd/internal/sqltest"
	"github.com/stretchr/testify/assert"
)

func Test_Company_GetAccounts(t *testing.T) {
	companyID := int64(42)
	companyID2 := int64(96)
	account1 := &Account{Account: &table.Account{CompanyID: &companyID}}
	account2 := &Account{Account: &table.Account{CompanyID: &companyID2}}
	company := &table.Company{ID: companyID}
	loadedSource := &companySource{accounts: []*Account{account1, account2}}
	tests := []struct {
		name     string
		company  *Company
		accounts []*Account
		err      error
	}{
		{"returns empty array for no accounts", &Company{Company: company, source: &companySource{accounts: []*Account{}}}, []*Account{}, nil},
		{"returns existing accounts", &Company{Company: company, source: loadedSource}, []*Account{account1}, nil},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.company.GetAccounts(nil)

			assert.Equal(t, test.accounts, result)
		})
	}
}

func Test_GetAllCompanies(t *testing.T) {
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		companies := []*table.Company{{ID: 1}}
		getAllCompaniesStub := mocka.Function(t, &getAllCompanies, companies)
		defer getAllCompaniesStub.Restore()

		result := GetAllCompanies(tx)

		assert.Equal(t, []interface{}{tx}, getAllCompaniesStub.GetFirstCall().Arguments())
		assert.Equal(t, companies[0], result[0].Company)
		assert.NotNil(t, result[0].source)
	})
}

func Test_GetCompanyByID(t *testing.T) {
	id := int64(42)
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		companies := []*table.Company{{ID: 1}}
		getCompanyByIDStub := mocka.Function(t, &getCompanyByID, companies)
		defer getCompanyByIDStub.Restore()

		result := GetCompanyByID(tx, id)

		assert.Equal(t, []interface{}{tx, id}, getCompanyByIDStub.GetFirstCall().Arguments())
		assert.Equal(t, companies[0], result[0].Company)
		assert.NotNil(t, result[0].source)
	})
}

func Test_GetCompanyByName(t *testing.T) {
	name := "the company"
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		companies := []*table.Company{{ID: 1}}
		getCompanyByNameStub := mocka.Function(t, &getCompanyByName, companies)
		defer getCompanyByNameStub.Restore()

		result := GetCompanyByName(tx, name)

		assert.Equal(t, []interface{}{tx, name}, getCompanyByNameStub.GetFirstCall().Arguments())
		assert.Equal(t, companies[0], result[0].Company)
		assert.NotNil(t, result[0].source)
	})
}

func Test_GetCompaniesByIDs(t *testing.T) {
	ids := []int64{42, 96}
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		companies := []*table.Company{{ID: 42}}
		getCompaniesByIDsStub := mocka.Function(t, &getCompaniesByIDs, companies)
		defer getCompaniesByIDsStub.Restore()

		result := GetCompaniesByIDs(tx, ids)

		assert.Equal(t, []interface{}{tx, ids}, getCompaniesByIDsStub.GetFirstCall().Arguments())
		assert.Equal(t, companies[0], result[0].Company)
		assert.NotNil(t, result[0].source)
	})
}

func Test_AddCompanies(t *testing.T) {
	id := int64(42)
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		company1 := &table.Company{ID: id}
		company2 := &table.Company{ID: id + 1}
		addCompanyStub := mocka.Function(t, &addCompany, company1)
		addCompanyStub.OnCall(1).Return(company2)
		defer addCompanyStub.Restore()
		validateNameStub := mocka.Function(t, &validateName)
		defer validateNameStub.Restore()

		results := AddCompanies(tx, []string{"company1", "company2"}, "somebody")

		assert.Equal(t, 2, validateNameStub.CallCount())
		assert.Equal(t, []interface{}{"company1"}, validateNameStub.GetCall(0).Arguments())
		assert.Equal(t, []interface{}{"company2"}, validateNameStub.GetCall(1).Arguments())
		assert.Equal(t, company1, results[0].Company)
		assert.NotNil(t, results[0].source)
		assert.Equal(t, company2, results[1].Company)
		assert.NotNil(t, results[1].source)
	})
}

func Test_UpdateCompanies(t *testing.T) {
	updates := []interface{}{
		map[string]interface{}{"id": 42, "name": "rename 42", "version": 1},
		map[string]interface{}{"id": 96, "name": "rename 96", "version": 1},
	}
	sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
		companies := []*Company{{Company: &table.Company{ID: 42, Name: "new name"}}}
		updateCompanyStub := mocka.Function(t, &updateCompany)
		defer updateCompanyStub.Restore()
		getCompaniesStub := mocka.Function(t, &GetCompaniesByIDs, companies)
		defer getCompaniesStub.Restore()
		validateNameStub := mocka.Function(t, &validateName)
		defer validateNameStub.Restore()

		result := UpdateCompanies(tx, updates, "somebody")

		assert.Equal(t, companies, result)
		assert.Equal(t, 2, updateCompanyStub.CallCount())
		assert.Equal(t, []interface{}{tx, int64(42), int64(1), "rename 42", "somebody"}, updateCompanyStub.GetCall(0).Arguments())
		assert.Equal(t, []interface{}{tx, int64(96), int64(1), "rename 96", "somebody"}, updateCompanyStub.GetCall(1).Arguments())
		assert.Equal(t, 1, getCompaniesStub.CallCount())
		assert.Equal(t, []interface{}{tx, []int64{42, 96}}, getCompaniesStub.GetCall(0).Arguments())
		assert.Equal(t, []interface{}{"rename 42"}, validateNameStub.GetCall(0).Arguments())
		assert.Equal(t, []interface{}{"rename 96"}, validateNameStub.GetCall(1).Arguments())
	})
}
