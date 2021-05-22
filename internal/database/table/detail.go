package table

// TransactionDetail represents a line item of a financial transaction.
type TransactionDetail struct {
	ID                    int64
	TransactionID         int64
	TransactionCategoryID *int64
	TransactionGroupID    *int64
	Memo                  *string
	Amount                float64
	AssetQuantity         *float64
	ExchangeAssetID       *int64
	RelatedDetailID       *int64
	Version               int
	Audited
}

func (d *TransactionDetail) PtrTo(column string) interface{} {
	switch column {
	case "id":
		return &d.ID
	case "transaction_id":
		return &d.TransactionID
	case "transaction_category_id":
		return &d.TransactionCategoryID
	case "transaction_group_id":
		return &d.TransactionGroupID
	case "memo":
		return &d.Memo
	case "amount":
		return &d.Amount
	case "asset_quantity":
		return &d.AssetQuantity
	case "exchange_asset_id":
		return &d.ExchangeAssetID
	case "related_detail_id":
		return &d.RelatedDetailID
	case "version":
		return &d.Version
	}
	return d.Audited.ptrToAudit(column)
}
