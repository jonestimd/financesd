package model

type TransactionDetail struct {
	ID              int
	TransactionID   int
	Transaction     *Transaction
	CategoryID      *int `gorm:"column:transaction_category_id"`
	GroupID         *int `gorm:"column:transaction_group_id"`
	Memo            *string
	Amount          float64
	AssetQuantity   float64
	ExchangeAssetID *int
	RelatedDetailID *int
	RelatedDetail   *TransactionDetail
	Version         int
	Audited
}
