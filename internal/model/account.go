package model

import (
	"database/sql"
	"encoding/json"
	"reflect"
)

// Account with a financial instutition.
type Account struct {
	ID               int64
	CompanyID        *int64
	Name             string
	Description      *string
	AccountNo        *string
	Type             string
	Closed           YesNo
	CurrencyID       int64
	Version          int64
	Balance          string
	TransactionCount int64
	source           *companySource
	Audited
}

func (a *Account) ptrTo(column string) interface{} {
	switch column {
	case "id":
		return &a.ID
	case "company_id":
		return &a.CompanyID
	case "name":
		return &a.Name
	case "description":
		return &a.Description
	case "account_no":
		return &a.AccountNo
	case "type":
		return &a.Type
	case "closed":
		return &a.Closed
	case "currency_id":
		return &a.CurrencyID
	case "version":
		return &a.Version
	case "balance":
		return &a.Balance
	case "transaction_count":
		return &a.TransactionCount
	}
	return a.Audited.ptrToAudit(column)
}

// GetCompany returns the company for the account.
func (a *Account) GetCompany(tx *sql.Tx) (*Company, error) {
	if a.CompanyID == nil {
		return nil, nil
	}
	if a.source.loadCompanies(tx) != nil {
		return nil, a.source.err
	}
	return a.source.companiesByID[*a.CompanyID], nil
}

var accountType = reflect.TypeOf(Account{})

const accountSQL = `select a.*,
	(select count(*) from transaction where account_id = a.id) transaction_count,
	(select sum(td.amount)
	 from transaction tx
	 join transaction_detail td on tx.id = td.transaction_id
	 left join transaction_category tc on td.transaction_category_id = tc.id
	 where tx.account_id = a.id and coalesce(tc.amount_type, '') != 'ASSET_VALUE') balance
from account a`

func runAccountQuery(tx *sql.Tx, cacheAccounts bool, query string, args ...interface{}) ([]*Account, error) {
	accounts, err := runQuery(tx, accountType, query, args...)
	if err != nil {
		return nil, err
	}
	return newCompanySource().setAccounts(accounts.([]*Account), cacheAccounts), nil
}

// GetAllAccounts loads all accounts.
func GetAllAccounts(tx *sql.Tx) ([]*Account, error) {
	return runAccountQuery(tx, true, accountSQL)
}

// GetAccountByID returns the account with ID.
func GetAccountByID(tx *sql.Tx, id int64) ([]*Account, error) {
	return runAccountQuery(tx, false, accountSQL+" where a.id = ?", id)
}

// GetAccountsByName returns the accounts having name.
func GetAccountsByName(tx *sql.Tx, name string) ([]*Account, error) {
	return runAccountQuery(tx, false, accountSQL+" where a.name = ?", name)
}

// getAccountsByCompanyIDs returns the accounts for the sepcified companies.
func getAccountsByCompanyIDs(tx *sql.Tx, companyIDs []int64) ([]*Account, error) {
	jsonIDs, _ := json.Marshal(companyIDs) // can't be cyclic, so ignoring error
	return runAccountQuery(tx, true, accountSQL+" where json_contains(?, cast(company_id as json))", jsonIDs)
}
