package database

import (
	"reflect"
	"time"
)

// Security represents an investment security.
type Security struct {
	AssetID          int64
	Type             string
	Shares           float64
	FirstAcquired    *time.Time
	CostBasis        *float64
	Dividends        *float64
	TransactionCount int64
	Asset
}

func (s *Security) ptrTo(column string) interface{} {
	switch column {
	case "asset_id":
		return &s.AssetID
	case "security_type":
		return &s.Type
	case "shares":
		return &s.Shares
	case "first_acquired":
		return &s.FirstAcquired
	case "cost_basis":
		return &s.CostBasis
	case "dividends":
		return &s.Dividends
	case "transaction_count":
		return &s.TransactionCount
	}
	return s.Asset.ptrTo(column)
}

var securityType = reflect.TypeOf(Security{})
