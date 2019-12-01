package graphql

import (
	"strconv"

	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	"github.com/jonestimd/financesd/internal/model"
)

// Schema
var securitySchema = graphql.NewObject(graphql.ObjectConfig{
	Name:        "security",
	Description: "an investement asset",
	Fields: addAsset(graphql.Fields{
		"type": &graphql.Field{Type: graphql.String},
	}),
})

var securityQueryFields = &graphql.Field{
	Type: graphql.NewList(securitySchema),
	Args: map[string]*graphql.ArgumentConfig{
		"id":     &graphql.ArgumentConfig{Type: graphql.ID, Description: "security ID"},
		"symbol": &graphql.ArgumentConfig{Type: graphql.String, Description: "unique security symbol"},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		db := p.Context.Value(DbContextKey).(*gorm.DB).Preload("Asset")
		securities := make([]*model.Security, 0)
		if id, ok := p.Args["id"]; ok {
			if intId, err := strconv.ParseInt(id.(string), 10, 64); err == nil {
				return securities, db.Find(&securities, "asset_id = ?", intId).Error
			} else {
				return nil, err
			}
		}
		if symbol, ok := p.Args["symbol"]; ok {
			return securities, db.Joins("join asset a on a.id = security.asset_id").
				Find(&securities, "a.symbol = ?", symbol.(string)).Error
		}
		return securities, db.Find(&securities).Error
	},
}
