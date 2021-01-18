package graphql

import (
	"database/sql"
	"strconv"

	"github.com/graphql-go/graphql"
)

// Schema
var securitySchema = graphql.NewObject(graphql.ObjectConfig{
	Name:        "security",
	Description: "an investement asset",
	Fields: addAudit(graphql.Fields{
		"id":         &graphql.Field{Type: graphql.ID},
		"assetType":  &graphql.Field{Type: graphql.String},
		"type":       &graphql.Field{Type: graphql.String},
		"name":       &graphql.Field{Type: graphql.String},
		"scale":      &graphql.Field{Type: graphql.Int},
		"symbol":     &graphql.Field{Type: graphql.String},
		"version":    &graphql.Field{Type: graphql.String},
		"changeUser": &graphql.Field{Type: graphql.String},
		"changeDate": &graphql.Field{Type: graphql.String},
	}),
})

const securitySQL = `select a.*, s.type security_type
from asset a
join security s on a.id = s.asset_id`

func securityColumnMapper(name string) string {
	switch name {
	case "security_type":
		return "type"
	case "type":
		return "asset_type"
	}
	return name
}

var securityQueryFields = &graphql.Field{
	Type: graphql.NewList(securitySchema),
	Args: map[string]*graphql.ArgumentConfig{
		"id":     {Type: graphql.ID, Description: "security ID"},
		"symbol": {Type: graphql.String, Description: "unique security symbol"},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		db := p.Context.Value(DbContextKey).(*sql.Tx)
		args := make([]interface{}, 0)
		sql := securitySQL
		if id, ok := p.Args["id"]; ok {
			intID, err := strconv.ParseInt(id.(string), 10, 64)
			if err != nil {
				return nil, err
			}
			args = append(args, intID)
			sql += " where a.id = ?"
		} else if symbol, ok := p.Args["symbol"]; ok {
			args = append(args, symbol)
			sql += " where a.symbol = ?"
		}
		rows, err := db.Query(sql, args...)
		if err != nil {
			return nil, err
		}
		return scanToMaps(rows, securityColumnMapper)
	},
}
