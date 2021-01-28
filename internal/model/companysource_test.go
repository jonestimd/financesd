package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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
			accounts := []*Account{
				{ID: 1, CompanyID: &companyID1},
				{ID: 2, CompanyID: &companyID1},
				{ID: 3},
				{ID: 4, CompanyID: &companyID2},
			}
			cs := newCompanySource()

			assert.Equal(t, accounts, cs.setAccounts(accounts, test.cacheAccounts))

			assert.ElementsMatch(t, []int64{companyID1, companyID2}, cs.companyIDs)
			if test.cacheAccounts {
				assert.Equal(t, accounts, cs.accounts)
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
	companies := []*Company{
		{ID: companyID1},
		{ID: companyID2},
	}
	cs := newCompanySource()

	assert.Equal(t, companies, cs.setCompanies(companies))

	assert.Equal(t, map[int64]*Company{
		companyID1: companies[0],
		companyID2: companies[1],
	}, cs.companiesByID)
	assert.ElementsMatch(t, []int64{companyID1, companyID2}, cs.companyIDs)
	for _, company := range companies {
		assert.Same(t, cs, company.source)
	}
}

func Test_companySource_loadCompanies(t *testing.T) {
	t.Run("returns existing error", func(t *testing.T) {
		cs := &companySource{err: errors.New("test error")}

		assert.Same(t, cs.err, cs.loadAccounts(nil))
	})
	t.Run("returns database error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			expectedErr := errors.New("test error")
			cs := newCompanySource()
			mockDB.ExpectQuery("select * from company where json_contains(?, cast(id as json))").
				WillReturnError(expectedErr)

			assert.Same(t, expectedErr, cs.loadCompanies(tx))

			assert.Same(t, expectedErr, cs.err)
		})
	})
	t.Run("loads companies", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			cs := &companySource{companyIDs: []int64{10, 20}}
			jsonIDs, _ := json.Marshal(cs.companyIDs)
			mockDB.ExpectQuery("select * from company where json_contains(?, cast(id as json))").
				WithArgs(jsonIDs).
				WillReturnRows(sqltest.MockRows("id").AddRow(cs.companyIDs[0]).AddRow(cs.companyIDs[1]))

			assert.Nil(t, cs.loadCompanies(tx))

			assert.Equal(t, 2, len(cs.companiesByID))
			assert.Nil(t, mockDB.ExpectationsWereMet())
			for _, id := range cs.companyIDs {
				assert.Equal(t, id, cs.companiesByID[id].ID)
				assert.Same(t, cs, cs.companiesByID[id].source)
			}
		})
	})
}

func Test_companySource_loadAccounts(t *testing.T) {
	t.Run("returns existing error", func(t *testing.T) {
		cs := &companySource{err: errors.New("test error")}

		assert.Same(t, cs.err, cs.loadAccounts(nil))
	})
	t.Run("returns database error", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			expectedErr := errors.New("test error")
			cs := newCompanySource()
			mockDB.ExpectQuery(accountSQL + " where json_contains(?, cast(company_id as json))").
				WillReturnError(expectedErr)

			assert.Same(t, expectedErr, cs.loadAccounts(tx))

			assert.Same(t, expectedErr, cs.err)
		})
	})
	t.Run("loads companies", func(t *testing.T) {
		sqltest.TestInTx(t, func(mockDB sqlmock.Sqlmock, tx *sql.Tx) {
			accountIDs := []int64{1, 2}
			cs := &companySource{companyIDs: []int64{10, 20}}
			jsonIDs, _ := json.Marshal(cs.companyIDs)
			mockDB.ExpectQuery(accountSQL + " where json_contains(?, cast(company_id as json))").
				WithArgs(jsonIDs).
				WillReturnRows(sqltest.MockRows("id").AddRow(accountIDs[0]).AddRow(accountIDs[1]))

			assert.Nil(t, cs.loadAccounts(tx))

			assert.Equal(t, 2, len(cs.accounts))
			assert.Nil(t, mockDB.ExpectationsWereMet())
			for i, account := range cs.accounts {
				assert.Equal(t, accountIDs[i], account.ID)
				assert.Same(t, cs, account.source)
			}
		})
	})
}
