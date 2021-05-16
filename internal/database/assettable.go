package database

// Asset represents a type of financial asset.
type Asset struct {
	ID      int64
	Name    string
	Type    string
	Scale   int
	Symbol  *string
	Version int
	Audited
}

func (a *Asset) ptrTo(column string) interface{} {
	switch column {
	case "id":
		return &a.ID
	case "name":
		return &a.Name
	case "type":
		return &a.Type
	case "scale":
		return &a.Scale
	case "symbol":
		return &a.Symbol
	case "version":
		return &a.Version
	}
	return a.Audited.ptrToAudit(column)
}
