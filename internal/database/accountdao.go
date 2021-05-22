package database

import (
	"database/sql"
	"encoding/json"
	"reflect"

	"github.com/jonestimd/financesd/internal/database/table"
)

var accountType = reflect.TypeOf(table.Account{})

const accountSQL = `select a.*,
	(select count(*) from transaction where account_id = a.id) transaction_count,
	(select sum(td.amount)
	 from transaction tx
	 join transaction_detail td on tx.id = td.transaction_id
	 left join transaction_category tc on td.transaction_category_id = tc.id
	 where tx.account_id = a.id and coalesce(tc.amount_type, '') != 'ASSET_VALUE') balance
from account a`

func runAccountQuery(tx *sql.Tx, query string, args ...interface{}) ([]*table.Account, error) {
	accounts, err := runQuery(tx, accountType, query, args...)
	if err != nil {
		return nil, err
	}
	return accounts.([]*table.Account), nil
}

// GetAllAccounts loads all accounts.
func GetAllAccounts(tx *sql.Tx) ([]*table.Account, error) {
	return runAccountQuery(tx, accountSQL)
}

// GetAccountByID returns the account with ID.
func GetAccountByID(tx *sql.Tx, id int64) ([]*table.Account, error) {
	return runAccountQuery(tx, accountSQL+" where a.id = ?", id)
}

// GetAccountsByName returns the accounts having name.
func GetAccountsByName(tx *sql.Tx, name string) ([]*table.Account, error) {
	return runAccountQuery(tx, accountSQL+" where a.name = ?", name)
}

// GetAccountsByCompanyIDs returns the accounts for the sepcified companies.
func GetAccountsByCompanyIDs(tx *sql.Tx, companyIDs []int64) ([]*table.Account, error) {
	jsonIDs, _ := json.Marshal(companyIDs) // can't be cyclic, so ignoring error
	return runAccountQuery(tx, accountSQL+" where json_contains(?, cast(company_id as json))", jsonIDs)
}
