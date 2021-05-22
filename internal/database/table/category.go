package table

// Category categorizes a transaction detail.
type Category struct {
	ID               int64
	Code             string
	Description      *string
	AmountType       string
	ParentID         *int64
	Security         *YesNo
	Income           *YesNo
	Version          int64
	TransactionCount int64
	Audited
}

func (tc *Category) PtrTo(column string) interface{} {
	switch column {
	case "id":
		return &tc.ID
	case "code":
		return &tc.Code
	case "description":
		return &tc.Description
	case "amount_type":
		return &tc.AmountType
	case "parent_id":
		return &tc.ParentID
	case "security":
		return &tc.Security
	case "income":
		return &tc.Income
	case "version":
		return &tc.Version
	case "transaction_count":
		return &tc.TransactionCount
	}
	return tc.Audited.ptrToAudit(column)
}
