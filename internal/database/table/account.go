package table

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

func (a *Account) PtrTo(column string) interface{} {
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
