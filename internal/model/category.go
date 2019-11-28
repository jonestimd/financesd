package model

type TransactionCategory struct {
	ID          int
	Code        string
	Description *string
	AmountType  string
	ParentId    *int
	Security    *YesNo
	Income      *YesNo
	Audited
}
