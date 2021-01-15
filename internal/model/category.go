package model

type TransactionCategory struct {
	ID          int64
	Code        string
	Description *string
	AmountType  string
	ParentId    *int64
	Security    *YesNo
	Income      *YesNo
	Audited
}
