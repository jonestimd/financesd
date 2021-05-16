package database

import (
	"database/sql"
)

const securitySummarySQL = `select t.security_id, count(distinct t.id) transaction_count
	, sum(coalesce(adjust_shares(t.security_id, t.date, td.asset_quantity), 0)) shares
	, min(t.date) first_acquired
	, sum(case when td.asset_quantity > 0
		then abs(td.amount)*(td.asset_quantity-coalesce(sl.purchase_shares,0))/td.asset_quantity
		else 0 end) cost_basis
	, sum(case when td.asset_quantity is null and td.amount > 0 then td.amount else 0 end) dividends
from tx t
join tx_detail td on t.id = td.tx_id
left join (
	select purchase_tx_detail_id, sum(purchase_shares) purchase_shares from security_lot
	group by purchase_tx_detail_id
) sl on td.id = sl.purchase_tx_detail_id
where t.security_id is not null`

const securitySQL = `select s.type security_type, a.*,
	coalesce(summary.transaction_count, 0) transaction_count,
	coalesce(summary.shares, 0) shares,
	summary.first_acquired, summary.cost_basis, summary.dividends
from security s
join asset a on s.asset_id = a.id
left join (
	` + securitySummarySQL + `
	group by t.security_id
) summary on summary.security_id = a.id`

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
