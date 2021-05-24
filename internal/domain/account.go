package domain

import (
	"database/sql"

	"github.com/graphql-go/graphql"
	"github.com/jonestimd/financesd/internal/database/table"
)

// Account with a financial instutition.
type Account struct {
	source *companySource
	*table.Account
}

func NewAccount(id int64, companyID *int64) *Account {
	return &Account{Account: &table.Account{ID: id, CompanyID: companyID}}
}

func (a *Account) Resolve(p graphql.ResolveParams) (interface{}, error) {
	return graphql.DefaultResolveFn(replaceSource(p, a.Account))
}

// GetCompany returns the company for the account.
func (a *Account) GetCompany(tx *sql.Tx) *Company {
	if a.CompanyID == nil {
		return nil
	}
	a.source.loadCompanies(tx)
	return a.source.companiesByID[*a.CompanyID]
}

// GetAllAccounts loads all accounts.
func GetAllAccounts(tx *sql.Tx) []*Account {
	accounts := getAllAccounts(tx)
	return newCompanySource().setAccounts(accounts, true)
}

// GetAccountByID returns the account with ID.
func GetAccountByID(tx *sql.Tx, id int64) []*Account {
	accounts := getAccountByID(tx, id)
	return newCompanySource().setAccounts(accounts, false)
}

// GetAccountsByName returns the accounts having name.
func GetAccountsByName(tx *sql.Tx, name string) []*Account {
	accounts := getAccountsByName(tx, name)
	return newCompanySource().setAccounts(accounts, false)
}

// GetAccountsByCompanyIDs returns the accounts for the sepcified companies.
func GetAccountsByCompanyIDs(tx *sql.Tx, companyIDs []int64) []*Account {
	accounts := getAccountsByCompanyIDs(tx, companyIDs)
	return newCompanySource().setAccounts(accounts, true)
}
