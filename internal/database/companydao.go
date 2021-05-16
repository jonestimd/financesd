package database

import (
	"database/sql"
	"fmt"
	"time"
)

func runCompanyQuery(tx *sql.Tx, query string, args ...interface{}) ([]*Company, error) {
	companies, err := runQuery(tx, companyType, query, args...)
	if err != nil {
		return nil, err
	}
	return companies.([]*Company), nil
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

// GetCompaniesByIDs loads specified companies.
func GetCompaniesByIDs(tx *sql.Tx, ids []int64) ([]*Company, error) {
	return runCompanyQuery(tx, "select * from company where json_contains(?, cast(id as json))", int64sToJson(ids))
}

// AddCompany adds a new company.
func AddCompany(tx *sql.Tx, name string, user string) (*Company, error) {
	changeDate := time.Now()
	id, err := runInsert(tx, "insert into company (name, change_user, change_date, version) values (?, ?, ?, 1)", name, user, changeDate)
	if err != nil {
		return nil, err
	}
	return &Company{ID: id, Name: name, Audited: Audited{user, &changeDate}, Version: 1}, nil
}

// DeleteCompanies deletes companies and returns the number of deleted companies.
func DeleteCompanies(tx *sql.Tx, ids []int) (int64, error) {
	return runUpdate(tx, "delete from company where json_contains(?, cast(id as json))", intsToJson(ids))
}

const updateCompanySQL = `update company set name = ?, change_date = current_timestamp, change_user = ?, version = version+1
where id = ? and version = ?`

// UpdateCompany updates a company name.
func UpdateCompany(tx *sql.Tx, id int64, version int64, name string, user string) error {
	if count, err := runUpdate(tx, updateCompanySQL, name, user, id, version); err != nil {
		return err
	} else if count == 0 {
		return fmt.Errorf("company not found (%d @ %d)", id, version)
	}
	return nil
}
