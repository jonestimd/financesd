package table

import (
	"time"
)

// Transaction represents a financial transaction.
type Transaction struct {
	ID              int64
	Date            time.Time
	Memo            *string
	ReferenceNumber *string
	Cleared         *YesNo
	AccountID       int64
	PayeeID         *int64
	SecurityID      *int64
	Details         []TransactionDetail
	Version         int
	Audited
}

func (t *Transaction) PtrTo(column string) interface{} {
	switch column {
	case "id":
		return &t.ID
	case "date":
		return &t.Date
	case "memo":
		return &t.Memo
	case "reference_number":
		return &t.ReferenceNumber
	case "cleared":
		return &t.Cleared
	case "account_id":
		return &t.AccountID
	case "payee_id":
		return &t.PayeeID
	case "security_id":
		return &t.SecurityID
	case "version":
		return &t.Version
	}
	return t.Audited.ptrToAudit(column)
}
