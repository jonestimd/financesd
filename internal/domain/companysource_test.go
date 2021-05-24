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

func Test_companySource(t *testing.T) {
	companyID1 := int64(10)
	companyID2 := int64(20)
	tests := []struct {
		name          string
		cacheAccounts bool
	}{
		{name: "sets source", cacheAccounts: false},
		{name: "sets accounts", cacheAccounts: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dbAccounts := []*table.Account{
				{ID: 1, CompanyID: &companyID1},
				{ID: 2, CompanyID: &companyID1},
				{ID: 3},
				{ID: 4, CompanyID: &companyID2},
			}
			cs := newCompanySource()

			accounts := cs.setAccounts(dbAccounts, test.cacheAccounts)

			assert.ElementsMatch(t, []int64{companyID1, companyID2}, cs.companyIDs)
			if test.cacheAccounts {
				assert.Equal(t, accounts, cs.accounts)
				for i, account := range accounts {
					assert.Same(t, dbAccounts[i], account.Account)
				}
			} else {
				assert.Nil(t, cs.accounts)
			}
			for _, account := range accounts {
				assert.Same(t, cs, account.source)
			}
		})
	}
}

func Test_companySource_setCompanies(t *testing.T) {
	companyID1 := int64(10)
	companyID2 := int64(20)
	dbCompanies := []*table.Company{
		{ID: companyID1},
		{ID: companyID2},
	}
	cs := newCompanySource()

	companies := cs.setCompanies(dbCompanies)

	assert.ElementsMatch(t, []int64{companyID1, companyID2}, cs.companyIDs)
	for i, company := range companies {
		assert.Same(t, dbCompanies[i], company.Company)
		assert.Same(t, cs, company.source)
		assert.Same(t, company, cs.companiesByID[company.ID])
	}
}

func Test_companySource_loadCompanies(t *testing.T) {
	sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
		cs := &companySource{companyIDs: []int64{10, 20}}
		mockDB.ExpectQuery("select * from company where json_contains(?, cast(id as json))").
			WithArgs("[10,20]").
			WillReturnRows(sqltest.MockRows("id").AddRow(cs.companyIDs[0]).AddRow(cs.companyIDs[1]))

		cs.loadCompanies(tx)

		assert.Equal(t, 2, len(cs.companiesByID))
		assert.Nil(t, mockDB.ExpectationsWereMet())
		for _, id := range cs.companyIDs {
			assert.Equal(t, id, cs.companiesByID[id].ID)
			assert.Same(t, cs, cs.companiesByID[id].source)
		}
	})
}

func Test_companySource_loadAccounts(t *testing.T) {
	sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
		accountIDs := []int64{1, 2}
		dbAccounts := []*table.Account{{ID: 1}, {ID: 2}}
		getAccountsStub := mocka.Function(t, &getAccountsByCompanyIDs, dbAccounts)
		defer getAccountsStub.Restore()
		cs := &companySource{companyIDs: []int64{10, 20}}

		cs.loadAccounts(tx)

		assert.Equal(t, 2, len(cs.accounts))
		for i, account := range cs.accounts {
			assert.Equal(t, accountIDs[i], account.ID)
			assert.Same(t, dbAccounts[i], account.Account)
			assert.Same(t, cs, account.source)
		}
	})
}
