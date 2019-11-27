package model

import "database/sql"

// Account with a financial instutition.
type Account struct {
	ID          int
	Company     *Company
	CompanyID   sql.NullInt64
	Name        string
	Description *string
	AccountNo   *string
	Type        string
	Closed      *YesNo
	CurrencyID  int
	Version     int
	Audited
}
