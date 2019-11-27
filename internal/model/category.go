package model

import "database/sql"

type TransactionCategory struct {
	ID          int
	Code        string
	Description *string
	AmountType  string
	ParentId    sql.NullInt64
	Security    *YesNo
	Income      *YesNo
	Audited
}
