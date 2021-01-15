package model

// Account with a financial instutition.
type Account struct {
	ID          int64
	Company     *Company
	CompanyID   *int64
	Name        string
	Description *string
	AccountNo   *string
	Type        string
	Closed      *YesNo
	CurrencyID  int64
	Version     int64
	Audited
}
