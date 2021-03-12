package model

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"time"
)

// Company contains information about a financial institution.
type Company struct {
	ID      int64
	Name    string
	Version int
	source  *companySource
	Audited
}

// PtrTo returns a pointer to the field for the database column.
func (c *Company) ptrTo(column string) interface{} {
	switch column {
	case "id":
		return &c.ID
	case "name":
		return &c.Name
	case "version":
		return &c.Version
	}
	return c.Audited.ptrToAudit(column)
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

var companyType = reflect.TypeOf(Company{})

func runCompanyQuery(tx *sql.Tx, query string, args ...interface{}) ([]*Company, error) {
	companies, err := runQuery(tx, companyType, query, args...)
	if err != nil {
		return nil, err
	}
	return newCompanySource().setCompanies(companies.([]*Company)), nil
}

// GetAllCompanies loads all companies.
func GetAllCompanies(tx *sql.Tx) ([]*Company, error) {
	return runCompanyQuery(tx, "select * from company")
}

// GetCompanyByID returns the company with the ID.
func GetCompanyByID(tx *sql.Tx, id int64) ([]*Company, error) {
	return runCompanyQuery(tx, "select * from company where id = ?", id)
}

// GetCompanyByName returns the company with the name.
func GetCompanyByName(tx *sql.Tx, name string) ([]*Company, error) {
	return runCompanyQuery(tx, "select * from company where name = ?", name)
}

// getCompaniesByIDs loads specified companies.
func getCompaniesByIDs(tx *sql.Tx, ids []int64) ([]*Company, error) {
	jsonIDs, _ := json.Marshal(ids) // can't be cyclic, so ignoring error
	return runCompanyQuery(tx, "select * from company where json_contains(?, cast(id as json))", jsonIDs)
}

// AddCompanies adds new companies.
func AddCompanies(tx *sql.Tx, names []string, user string) ([]*Company, error) {
	changeDate := time.Now()
	companies := make([]*Company, len(names))
	for i, name := range names {
		if err := validateName(name); err != nil {
			return nil, err
		}
		id, err := runInsert(tx, "insert into company (name, change_user, change_date, version) values (?, ?, ?, 1)", name, user, changeDate)
		if err != nil {
			return nil, err
		}
		companies[i] = &Company{ID: id, Name: name, Audited: Audited{user, &changeDate}, Version: 1}
	}
	return newCompanySource().setCompanies(companies), nil
}

// DeleteCompanies deletes companies.
func DeleteCompanies(tx *sql.Tx, ids []int, user string) (int64, error) {
	rs, err := runUpdate(tx, "delete from company where json_contains(?, cast(id as json))", intsToJson(ids))
	if err != nil {
		return 0, err
	}
	return rs.RowsAffected()
}
