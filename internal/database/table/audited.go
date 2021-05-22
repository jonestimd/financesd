package table

import "time"

// Audited contains audit information for a database record.
type Audited struct {
	ChangeUser string
	ChangeDate *time.Time
}

func (a *Audited) ptrToAudit(column string) interface{} {
	switch column {
	case "change_user":
		return &a.ChangeUser
	case "change_date":
		return &a.ChangeDate
	}
	return nil
}
