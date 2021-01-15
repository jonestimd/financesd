package model

type TransactionDetail struct {
	ID                    int64
	TransactionID         int64
	Transaction           *Transaction
	TransactionCategoryID *int64
	TransactionGroupID    *int64
	Memo                  *string
	Amount                float64
	AssetQuantity         float64
	ExchangeAssetID       *int64
	RelatedDetailID       *int64
	RelatedDetail         *TransactionDetail
	Version               int
	Audited
}
