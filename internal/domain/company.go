package domain

import (
	"database/sql"

	"github.com/graphql-go/graphql"
	"github.com/jonestimd/financesd/internal/database/table"
)

// Company contains information about a financial institution.
type Company struct {
	source *companySource
	*table.Company
}

func NewCompany(id int64, name string) *Company {
	return &Company{Company: &table.Company{ID: id, Name: name}}
}

func (c *Company) Resolve(p graphql.ResolveParams) (interface{}, error) {
	return graphql.DefaultResolveFn(replaceSource(p, c.Company))
}

// GetAccounts returns the accounts for the company.
func (c *Company) GetAccounts(tx *sql.Tx) []*Account {
	if c.source.accounts == nil {
		c.source.loadAccounts(tx)
	}
	accounts := make([]*Account, 0, 1)
	for _, account := range c.source.accounts {
		if *account.CompanyID == c.ID {
			accounts = append(accounts, account)
		}
	}
	return accounts
}

// GetAllCompanies loads all companies.
func GetAllCompanies(tx *sql.Tx) []*Company {
	return newCompanySource().setCompanies(getAllCompanies(tx))
}

// GetCompanyByID returns the company with the ID.
func GetCompanyByID(tx *sql.Tx, id int64) []*Company {
	return newCompanySource().setCompanies(getCompanyByID(tx, id))
}

// GetCompanyByName returns the company with the name.
func GetCompanyByName(tx *sql.Tx, name string) []*Company {
	return newCompanySource().setCompanies(getCompanyByName(tx, name))
}

// GetCompaniesByIDs loads specified companies.
var GetCompaniesByIDs = func(tx *sql.Tx, ids []int64) []*Company {
	return newCompanySource().setCompanies(getCompaniesByIDs(tx, ids))
}

// AddCompanies adds new companies.
func AddCompanies(tx *sql.Tx, names []string, user string) []*Company {
	companies := make([]*table.Company, len(names))
	for i, name := range names {
		validateName(name)
		companies[i] = addCompany(tx, name, user)
	}
	return newCompanySource().setCompanies(companies)
}

// UpdateCompanies updates company names.
func UpdateCompanies(tx *sql.Tx, args interface{}, user string) []*Company {
	updates := args.([]interface{})
	ids := make([]int64, len(updates))
	for i, company := range updates {
		values := company.(map[string]interface{})
		ids[i] = int64(values["id"].(int))
		name := values["name"].(string)
		validateName(name)
		version := int64(values["version"].(int))
		updateCompany(tx, ids[i], version, name, user)
	}
	return GetCompaniesByIDs(tx, ids)
}
