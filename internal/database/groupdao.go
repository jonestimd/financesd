package database

import (
	"database/sql"
	"reflect"

	"github.com/jonestimd/financesd/internal/database/table"
)

var groupType = reflect.TypeOf(table.Group{})

const groupSQL = `select g.*,
	(select count(distinct transaction_id) from transaction_detail where transaction_group_id = g.id) transaction_count
from transaction_group g`

// GetAllGroups loads all groups.
func GetAllGroups(tx *sql.Tx) ([]*table.Group, error) {
	groups, err := runQuery(tx, groupType, groupSQL)
	if err != nil {
		return nil, err
	}
	return groups.([]*table.Group), nil
}
