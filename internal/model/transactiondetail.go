package model

type TransactionDetail struct {
	ID                    int
	TransactionID         int
	Transaction           *Transaction
	TransactionCategoryID *int
	TransactionGroupID    *int
	Memo                  *string
	Amount                float64
	AssetQuantity         float64
	ExchangeAssetID       *int
	RelatedDetailID       *int
	RelatedDetail         *TransactionDetail
	Version               int
	Audited
}
