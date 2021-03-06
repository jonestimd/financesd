package database

import (
	"database/sql"
	"reflect"

	"github.com/jonestimd/financesd/internal/database/table"
)

var categoryType = reflect.TypeOf(table.Category{})

const categorySQL = `select c.*,
	(select count(distinct transaction_id) from transaction_detail where transaction_category_id = c.id) transaction_count
from transaction_category c`

// GetAllCategories loads all transaction categories.
func GetAllCategories(tx *sql.Tx) []*table.Category {
	categories := runQuery(tx, categoryType, categorySQL)
	return categories.([]*table.Category)
}
