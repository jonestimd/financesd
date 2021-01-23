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

var accountType = reflect.TypeOf(Account{})

// Accounts contains result of loading accounts.
type Accounts struct {
	Accounts   []*Account
	err        error
	companyIDs []int64
	byID       map[int64]*Account
}

// NewAccounts creates an error result.
func NewAccounts(err error) *Accounts {
	return &Accounts{err: err}
}

// CompanyIDs returns the unique company IDs for the accounts.
func (a *Accounts) CompanyIDs() []int64 {
	if a.companyIDs == nil {
		companyIDs := newIDSet()
		a.byID = make(map[int64]*Account, len(a.Accounts))
		for _, account := range a.Accounts {
			a.byID[account.ID] = account
			if account.CompanyID != nil {
				companyIDs.Add(*account.CompanyID)
			}
		}
		a.companyIDs = companyIDs.Values()
	}
	return a.companyIDs
}

// ByID returns a map of accounts keyed by ID.
func (a *Accounts) ByID() map[int64]*Account {
	if a.companyIDs == nil {
		a.CompanyIDs()
	}
	return a.byID
}

// Result returns the accounts as an interface.
func (a *Accounts) Result() interface{} {
	return a.Accounts
}

func (a *Accounts) Error() error {
	return a.err
}

const accountSQL = `select a.*,
	(select count(*) from transaction where account_id = a.id) transaction_count,
	(select sum(td.amount)
	 from transaction tx
	 join transaction_detail td on tx.id = td.transaction_id
	 left join transaction_category tc on td.transaction_category_id = tc.id
	 where tx.account_id = a.id and coalesce(tc.amount_type, '') != 'ASSET_VALUE') balance
from account a`

// GetAllAccounts loads all accounts.
func GetAllAccounts(tx *sql.Tx) *Accounts {
	accounts, err := runQuery(tx, accountType, accountSQL)
	if err != nil {
		return &Accounts{err: err}
	}
	return &Accounts{Accounts: accounts.([]*Account)}
}

// GetAccountByID returns the account with ID.
func GetAccountByID(tx *sql.Tx, id int64) *Accounts {
	accounts, err := runQuery(tx, accountType, accountSQL+" where a.id = ?", id)
	if err != nil {
		return &Accounts{err: err}
	}
	return &Accounts{Accounts: accounts.([]*Account)}
}

// GetAccountsByName returns the accounts having name.
func GetAccountsByName(tx *sql.Tx, name string) *Accounts {
	accounts, err := runQuery(tx, accountType, accountSQL+" where a.name = ?", name)
	if err != nil {
		return &Accounts{err: err}
	}
	return &Accounts{Accounts: accounts.([]*Account)}
}

// GetAccountsByCompanyIDs returns the accounts for the sepcified companies.
func GetAccountsByCompanyIDs(tx *sql.Tx, companyIDs []int64) *Accounts {
	jsonIDs, _ := json.Marshal(companyIDs) // can't be cyclic, so ignoring error
	accounts, err := runQuery(tx, accountType, accountSQL+" where json_contains(?, cast(company_id as json))", jsonIDs)
	if err != nil {
		return &Accounts{err: err}
	}
	return &Accounts{Accounts: accounts.([]*Account)}
}
