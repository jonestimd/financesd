package model

// Account with a financial instutition.
type Account struct {
	ID          int
	Company     *Company
	CompanyID   *int
	Name        string
	Description *string
	AccountNo   *string
	Type        string
	Closed      *YesNo
	CurrencyID  int
	Version     int
	Audited
}
