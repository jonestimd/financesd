package model

import "database/sql"

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

func (cs *companySource) setAccounts(accounts []*Account, cacheAccounts bool) []*Account {
	if cacheAccounts {
		cs.accounts = accounts
	}
	companyIDs := newIDSet()
	for _, account := range accounts {
		account.source = cs
		if account.CompanyID != nil {
			companyIDs.Add(*account.CompanyID)
		}
	}
	cs.companyIDs = companyIDs.Values()
	return accounts
}

func (cs *companySource) setCompanies(companies []*Company) []*Company {
	cs.companiesByID = make(map[int64]*Company, len(companies))
	cs.companyIDs = make([]int64, len(companies))
	for i, company := range companies {
		company.source = cs
		cs.companyIDs[i] = company.ID
		cs.companiesByID[company.ID] = company
	}
	return companies
}

func (cs *companySource) loadCompanies(tx *sql.Tx) error {
	if cs.err == nil && cs.companiesByID == nil {
		if companies, err := getCompaniesByIDs(tx, cs.companyIDs); err != nil {
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
		if accounts, err := getAccountsByCompanyIDs(tx, cs.companyIDs); err != nil {
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
