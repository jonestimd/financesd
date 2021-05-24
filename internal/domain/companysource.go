package domain

import (
	"database/sql"

	"github.com/jonestimd/financesd/internal/database/table"
)

// companySource provides companies and accounts for a GraphQL request.
type companySource struct {
	accounts      []*Account
	companyIDs    []int64
	companiesByID map[int64]*Company
}

func newCompanySource() *companySource {
	return &companySource{}
}

func (cs *companySource) setAccounts(accounts []*table.Account, cacheAccounts bool) []*Account {
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

func (cs *companySource) setCompanies(dbCompanies []*table.Company) []*Company {
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

func (cs *companySource) loadCompanies(tx *sql.Tx) {
	if cs.companiesByID == nil {
		companies := GetCompaniesByIDs(tx, cs.companyIDs)
		cs.companiesByID = make(map[int64]*Company, len(companies))
		for _, company := range companies {
			company.source = cs
			cs.companiesByID[company.ID] = company
		}
	}
}

func (cs *companySource) loadAccounts(tx *sql.Tx) {
	if cs.accounts == nil {
		accounts := GetAccountsByCompanyIDs(tx, cs.companyIDs)
		cs.accounts = accounts
		for _, account := range accounts {
			account.source = cs
		}
	}
}
