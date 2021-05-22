package table

// Group represents an alternate categorization for a transaction detail.
type Group struct {
	ID               int64
	Name             string
	Description      *string
	Version          int
	TransactionCount int64
	Audited
}

func (g *Group) PtrTo(column string) interface{} {
	switch column {
	case "id":
		return &g.ID
	case "name":
		return &g.Name
	case "description":
		return &g.Description
	case "version":
		return &g.Version
	case "transaction_count":
		return &g.TransactionCount
	}
	return g.Audited.ptrToAudit(column)
}
