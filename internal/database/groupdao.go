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
func GetAllGroups(tx *sql.Tx) []*table.Group {
	groups := runQuery(tx, groupType, groupSQL)
	return groups.([]*table.Group)
}
