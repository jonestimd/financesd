package model

import (
	"database/sql"
	"reflect"
)

// Group represents an alternate categorization for a transaction detail.
type Group struct {
	ID               int64
	Name             string
	Description      *string
	Version          int
	TransactionCount int64
	Audited
}

func (g *Group) ptrTo(column string) interface{} {
	switch column {
	case "id":
		return &g.ID
	case "name":
		return &g.Name
	case "description":
		return &g.Description
	case "version":
		return &g.Version
	case "transaction_count":
		return &g.TransactionCount
	}
	return g.Audited.ptrToAudit(column)
}

var groupType = reflect.TypeOf(Group{})

const groupSQL = `select g.*,
	(select count(distinct transaction_id) from transaction_detail where transaction_group_id = g.id) transaction_count
from transaction_group g`

// GetAllGroups loads all groups.
func GetAllGroups(tx *sql.Tx) ([]*Group, error) {
	groups, err := runQuery(tx, groupType, groupSQL)
	if err != nil {
		return nil, err
	}
	return groups.([]*Group), nil
}
