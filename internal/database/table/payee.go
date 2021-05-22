package table

// Payee represents the other party party in a financial transaction.
type Payee struct {
	ID               int64
	Name             string
	Version          int
	TransactionCount int64
	Audited
}

func (p *Payee) PtrTo(column string) interface{} {
	switch column {
	case "id":
		return &p.ID
	case "name":
		return &p.Name
	case "version":
		return &p.Version
	case "transaction_count":
		return &p.TransactionCount
	}
	return p.Audited.ptrToAudit(column)
}
