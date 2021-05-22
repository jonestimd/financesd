package database

import (
	"database/sql"
	"reflect"

	"github.com/jonestimd/financesd/internal/database/table"
)

var payeeType = reflect.TypeOf(table.Payee{})

const payeeSQL = `select p.*,
	(select count(t.id) from transaction t where t.payee_id = p.id) transaction_count
from payee p`

// GetAllPayees loads all payees.
func GetAllPayees(tx *sql.Tx) ([]*table.Payee, error) {
	payees, err := runQuery(tx, payeeType, payeeSQL)
	if err != nil {
		return nil, err
	}
	return payees.([]*table.Payee), nil
}
