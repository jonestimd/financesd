package model

import (
	"database/sql"

	"github.com/graphql-go/graphql"
	"github.com/jonestimd/financesd/internal/database"
)

// Company contains information about a financial institution.
type Company struct {
	source *companySource
	*database.Company
}

func NewCompany(id int64, name string) *Company {
	return &Company{Company: &database.Company{ID: id, Name: name}}
}

func (c *Company) Resolve(p graphql.ResolveParams) (interface{}, error) {
	return graphql.DefaultResolveFn(replaceSource(p, c.Company))
}

// GetAccounts returns the accounts for the company.
func (c *Company) GetAccounts(tx *sql.Tx) ([]*Account, error) {
	if c.source.accounts == nil && c.source.loadAccounts(tx) != nil {
		return nil, c.source.err
	}
	accounts := make([]*Account, 0, 1)
	for _, account := range c.source.accounts {
		if *account.CompanyID == c.ID {
			accounts = append(accounts, account)
		}
	}
	return accounts, nil
}

func companiesOrError(companies []*database.Company, err error) ([]*Company, error) {
	if err != nil {
		return nil, err
	}
	return newCompanySource().setCompanies(companies), nil
}

// GetAllCompanies loads all companies.
func GetAllCompanies(tx *sql.Tx) ([]*Company, error) {
	return companiesOrError(getAllCompanies(tx))
}

// GetCompanyByID returns the company with the ID.
func GetCompanyByID(tx *sql.Tx, id int64) ([]*Company, error) {
	return companiesOrError(getCompanyByID(tx, id))
}

// GetCompanyByName returns the company with the name.
func GetCompanyByName(tx *sql.Tx, name string) ([]*Company, error) {
	return companiesOrError(getCompanyByName(tx, name))
}

// GetCompaniesByIDs loads specified companies.
var GetCompaniesByIDs = func(tx *sql.Tx, ids []int64) ([]*Company, error) {
	return companiesOrError(getCompaniesByIDs(tx, ids))
}

// AddCompanies adds new companies.
func AddCompanies(tx *sql.Tx, names []string, user string) ([]*Company, error) {
	companies := make([]*database.Company, len(names))
	var err error
	for i, name := range names {
		if err = validateName(name); err != nil {
			return nil, err
		}
		if companies[i], err = addCompany(tx, name, user); err != nil {
			return nil, err
		}
	}
	return newCompanySource().setCompanies(companies), nil
}

// UpdateCompanies updates company names.
func UpdateCompanies(tx *sql.Tx, args interface{}, user string) ([]*Company, error) {
	updates := args.([]interface{})
	ids := make([]int64, len(updates))
	for i, company := range updates {
		values := company.(map[string]interface{})
		ids[i] = int64(values["id"].(int))
		name := values["name"].(string)
		version := int64(values["version"].(int))
		if err := updateCompany(tx, ids[i], version, name, user); err != nil {
			return nil, err
		}
	}
	return GetCompaniesByIDs(tx, ids)
}
