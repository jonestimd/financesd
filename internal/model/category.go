package model

import (
	"database/sql"
	"reflect"
)

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

func (tc *Category) ptrTo(column string) interface{} {
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

var categoryType = reflect.TypeOf(Category{})

const categorySQL = `select c.*,
	(select count(distinct transaction_id) from transaction_detail where transaction_category_id = c.id) transaction_count
from transaction_category c`

// GetAllCategories loads all transaction categories.
func GetAllCategories(tx *sql.Tx) ([]*Category, error) {
	categories, err := runQuery(tx, categoryType, categorySQL)
	if err != nil {
		return nil, err
	}
	return categories.([]*Category), nil
}
