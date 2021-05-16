package model

import (
	"database/sql"

	"github.com/jonestimd/financesd/internal/database"
)

// companySource provides companies and accounts for a GraphQL request.
type companySource struct {
	accounts      []*Account
	companyIDs    []int64
	companiesByID map[int64]*Company
	err           error
}

func newCompanySource() *companySource {
	return &companySource{}
}

func (cs *companySource) setAccounts(accounts []*database.Account, cacheAccounts bool) []*Account {
	companyIDs := newIDSet()
	accts := make([]*Account, len(accounts))
	for i, account := range accounts {
		accts[i] = &Account{source: cs, Account: account}
		if account.CompanyID != nil {
			companyIDs.Add(*account.CompanyID)
		}
	}
	cs.companyIDs = companyIDs.Values()
	if cacheAccounts {
		cs.accounts = accts
	}
	return accts
}

func (cs *companySource) setCompanies(dbCompanies []*database.Company) []*Company {
	cs.companiesByID = make(map[int64]*Company, len(dbCompanies))
	cs.companyIDs = make([]int64, len(dbCompanies))
	companies := make([]*Company, len(dbCompanies))
	for i, company := range dbCompanies {
		companies[i] = &Company{source: cs, Company: company}
		cs.companyIDs[i] = company.ID
		cs.companiesByID[company.ID] = companies[i]
	}
	return companies
}

func (cs *companySource) loadCompanies(tx *sql.Tx) error {
	if cs.err == nil && cs.companiesByID == nil {
		if companies, err := GetCompaniesByIDs(tx, cs.companyIDs); err != nil {
			cs.err = err
		} else {
			cs.companiesByID = make(map[int64]*Company, len(companies))
			for _, company := range companies {
				company.source = cs
				cs.companiesByID[company.ID] = company
			}
		}
	}
	return cs.err
}

func (cs *companySource) loadAccounts(tx *sql.Tx) error {
	if cs.err == nil && cs.accounts == nil {
		if accounts, err := GetAccountsByCompanyIDs(tx, cs.companyIDs); err != nil {
			cs.err = err
		} else {
			cs.accounts = accounts
			for _, account := range accounts {
				account.source = cs
			}
		}
	}
	return cs.err
}
