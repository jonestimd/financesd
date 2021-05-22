package table

// Company contains information about a financial institution.
type Company struct {
	ID      int64
	Name    string
	Version int
	Audited
}

// PtrTo returns a pointer to the field for the database column.
func (c *Company) PtrTo(column string) interface{} {
	switch column {
	case "id":
		return &c.ID
	case "name":
		return &c.Name
	case "version":
		return &c.Version
	}
	return c.Audited.ptrToAudit(column)
}
