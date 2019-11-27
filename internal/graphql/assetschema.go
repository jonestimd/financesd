package graphql

import (
	"strconv"

	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	"github.com/jonestimd/financesd/internal/model"
)

// Schema
var assetSchema = graphql.NewObject(graphql.ObjectConfig{
	Name:        "asset",
	Description: "a type of asset",
	Fields: addAudit(graphql.Fields{
		"id":     &graphql.Field{Type: graphql.ID},
		"name":   &graphql.Field{Type: graphql.String},
		"type":   &graphql.Field{Type: graphql.String},
		"scale":  &graphql.Field{Type: graphql.Int},
		"symbol": &graphql.Field{Type: graphql.String},
	}),
})

var assetQueryFields = &graphql.Field{
	Type: graphql.NewList(assetSchema),
	Args: map[string]*graphql.ArgumentConfig{
		"id":   &graphql.ArgumentConfig{Type: graphql.ID, Description: "asset ID"},
		"name": &graphql.ArgumentConfig{Type: graphql.String, Description: "unique asset name"},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		db := p.Context.Value(DbContextKey).(*gorm.DB)
		assets := make([]*model.Asset, 0)
		if id, ok := p.Args["id"]; ok {
			if intId, err := strconv.ParseInt(id.(string), 10, 64); err == nil {
				return assets, db.Find(&assets, intId).Error
			} else {
				return nil, err
			}
		}
		if name, ok := p.Args["name"]; ok {
			return assets, db.Find(&assets, "lower(name) = lower(?)", name.(string)).Error
		}
		return assets, db.Order("name").Find(&assets).Error
	},
}
