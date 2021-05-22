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
func (a *Account) GetCompany(tx *sql.Tx) (*Company, error) {
	if a.CompanyID == nil {
		return nil, nil
	}
	if a.source.loadCompanies(tx) != nil {
		return nil, a.source.err
	}
	return a.source.companiesByID[*a.CompanyID], nil
}

func accountsOrError(accounts []*table.Account, err error, cacheAccounts bool) ([]*Account, error) {
	if err != nil {
		return nil, err
	}
	return newCompanySource().setAccounts(accounts, cacheAccounts), nil
}

// GetAllAccounts loads all accounts.
func GetAllAccounts(tx *sql.Tx) ([]*Account, error) {
	accounts, err := getAllAccounts(tx)
	return accountsOrError(accounts, err, true)
}

// GetAccountByID returns the account with ID.
func GetAccountByID(tx *sql.Tx, id int64) ([]*Account, error) {
	accounts, err := getAccountByID(tx, id)
	return accountsOrError(accounts, err, false)
}

// GetAccountsByName returns the accounts having name.
func GetAccountsByName(tx *sql.Tx, name string) ([]*Account, error) {
	accounts, err := getAccountsByName(tx, name)
	return accountsOrError(accounts, err, false)
}

// GetAccountsByCompanyIDs returns the accounts for the sepcified companies.
func GetAccountsByCompanyIDs(tx *sql.Tx, companyIDs []int64) ([]*Account, error) {
	accounts, err := getAccountsByCompanyIDs(tx, companyIDs)
	return accountsOrError(accounts, err, true)
}
