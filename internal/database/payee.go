package database

import (
	"database/sql"
	"reflect"
)

// Payee represents the other party party in a financial transaction.
type Payee struct {
	ID               int64
	Name             string
	Version          int
	TransactionCount int64
	Audited
}

func (p *Payee) ptrTo(column string) interface{} {
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

var payeeType = reflect.TypeOf(Payee{})

const payeeSQL = `select p.*,
	(select count(t.id) from transaction t where t.payee_id = p.id) transaction_count
from payee p`

// GetAllPayees loads all payees.
func GetAllPayees(tx *sql.Tx) ([]*Payee, error) {
	payees, err := runQuery(tx, payeeType, payeeSQL)
	if err != nil {
		return nil, err
	}
	return payees.([]*Payee), nil
}
