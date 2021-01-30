package schema

import (
	"database/sql"
	"strconv"

	"github.com/graphql-go/graphql"
)

// Schema
var securitySchema = graphql.NewObject(graphql.ObjectConfig{
	Name:        "security",
	Description: "an investement asset",
	Fields: graphql.Fields{
		"id":               &graphql.Field{Type: graphql.ID, Resolve: nestedResolver("Asset", "ID")},
		"assetType":        &graphql.Field{Type: graphql.String, Resolve: nestedResolver("Asset", "Type")},
		"type":             &graphql.Field{Type: graphql.String},
		"name":             &graphql.Field{Type: graphql.String, Resolve: nestedResolver("Asset", "Name")},
		"scale":            &graphql.Field{Type: graphql.Int, Resolve: nestedResolver("Asset", "Scale")},
		"symbol":           &graphql.Field{Type: graphql.String, Resolve: nestedResolver("Asset", "Symbol")},
		"version":          &graphql.Field{Type: graphql.String, Resolve: nestedResolver("Asset", "Version")},
		"transactionCount": &graphql.Field{Type: graphql.Int},
		"changeUser":       &graphql.Field{Type: graphql.String, Resolve: nestedResolver("Asset", "ChangeUser")},
		"changeDate":       &graphql.Field{Type: graphql.String, Resolve: nestedResolver("Asset", "ChangeDate")},
	},
})

var securityQueryFields = &graphql.Field{
	Type: graphql.NewList(securitySchema),
	Args: map[string]*graphql.ArgumentConfig{
		"id":     {Type: graphql.ID, Description: "security ID"},
		"symbol": {Type: graphql.String, Description: "unique security symbol"},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		tx := p.Context.Value(DbContextKey).(*sql.Tx)
		if id, ok := p.Args["id"]; ok {
			intID, err := strconv.ParseInt(id.(string), 10, 64)
			if err != nil {
				return nil, err
			}
			return getSecurityByID(tx, intID)
		} else if symbol, ok := p.Args["symbol"]; ok {
			return getSecurityBySymbol(tx, symbol.(string))
		}
		return getAllSecurities(tx)
	},
}
