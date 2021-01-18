package model

import (
	"database/sql"
	"encoding/json"
	"reflect"
)

// Company contains information about a financial institution.
type Company struct {
	ID      int64
	Name    string
	Version int
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

var companyType = reflect.TypeOf(Company{})

// Companies contains result of loading companies.
type Companies struct {
	Companies []*Company
	err       error
}

// CompanyIDs returns the IDs of the companies.
func (c *Companies) CompanyIDs() []int64 {
	companyIDs := make([]int64, len(c.Companies))
	for i, company := range c.Companies {
		companyIDs[i] = (company.ID)
	}
	return companyIDs
}

// Result returns the companies as an interface.
func (c *Companies) Result() interface{} {
	return c.Companies
}

func (c *Companies) Error() error {
	return c.err
}

// GetAllCompanies loads all companies.
func GetAllCompanies(tx *sql.Tx) *Companies {
	companies, err := runQuery(tx, companyType, "select * from company")
	if err != nil {
		return &Companies{err: err}
	}
	return &Companies{Companies: companies.([]*Company)}
}

// GetCompanyByID returns the company with the ID.
func GetCompanyByID(tx *sql.Tx, id int64) *Companies {
	companies, err := runQuery(tx, companyType, "select * from company where id = ?", id)
	if err != nil {
		return &Companies{err: err}
	}
	return &Companies{Companies: companies.([]*Company)}
}

// GetCompanyByName returns the company with the name.
func GetCompanyByName(tx *sql.Tx, name string) *Companies {
	companies, err := runQuery(tx, companyType, "select * from company where name = ?", name)
	if err != nil {
		return &Companies{err: err}
	}
	return &Companies{Companies: companies.([]*Company)}
}

// GetCompaniesByIDs loads specified companies.
func GetCompaniesByIDs(tx *sql.Tx, ids []int64) (map[int64]*Company, error) {
	jsonIDs, _ := json.Marshal(ids) // can't be cyclic, so ignoring error
	models, err := runQuery(tx, companyType, "select * from company where json_contains(?, cast(id as json))", jsonIDs)
	if err != nil {
		return nil, err
	}
	companies := models.([]*Company)
	byID := make(map[int64]*Company, len(companies))
	for _, company := range companies {
		byID[company.ID] = company
	}
	return byID, nil
}
