package model

import (
	"database/sql"
	"reflect"
)

// Security represents an investment security.
type Security struct {
	AssetID          int64
	Type             string
	TransactionCount int64
	Asset
}

func (s *Security) ptrTo(column string) interface{} {
	switch column {
	case "asset_id":
		return &s.AssetID
	case "security_type":
		return &s.Type
	case "transaction_count":
		return &s.TransactionCount
	}
	return s.Asset.ptrTo(column)
}

var securityType = reflect.TypeOf(Security{})

const securitySQL = `select s.type security_type, a.*,
	(select count(t.id) from transaction t where t.security_id = a.id) transaction_count
from security s
join asset a on s.asset_id = a.id`

// GetAllSecurities loads all securities.
func GetAllSecurities(tx *sql.Tx) ([]*Security, error) {
	securities, err := runQuery(tx, securityType, securitySQL)
	if err != nil {
		return nil, err
	}
	return securities.([]*Security), nil
}

// GetSecurityByID returns the security with ID.
func GetSecurityByID(tx *sql.Tx, id int64) ([]*Security, error) {
	securities, err := runQuery(tx, securityType, securitySQL+" where a.id = ?", id)
	if err != nil {
		return nil, err
	}
	return securities.([]*Security), nil
}

// GetSecurityBySymbol returns the security for the symbol.
func GetSecurityBySymbol(tx *sql.Tx, symbol string) ([]*Security, error) {
	securities, err := runQuery(tx, securityType, securitySQL+" where s.symbol = ?", symbol)
	if err != nil {
		return nil, err
	}
	return securities.([]*Security), nil
}
