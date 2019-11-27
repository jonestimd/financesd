package graphql

import (
	"strconv"

	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	"github.com/jonestimd/financesd/internal/model"
)

// Schema
var payeeSchema = graphql.NewObject(graphql.ObjectConfig{
	Name:        "payee",
	Description: "the other party in a transaction",
	Fields: addAudit(graphql.Fields{
		"id":   &graphql.Field{Type: graphql.ID},
		"name": &graphql.Field{Type: graphql.String},
	}),
})

var payeeQueryFields = &graphql.Field{
	Type: graphql.NewList(payeeSchema),
	Args: map[string]*graphql.ArgumentConfig{
		"id":   &graphql.ArgumentConfig{Type: graphql.ID, Description: "payee ID"},
		"name": &graphql.ArgumentConfig{Type: graphql.String, Description: "unique payee name"},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		db := p.Context.Value(DbContextKey).(*gorm.DB)
		payees := make([]*model.Payee, 0)
		if id, ok := p.Args["id"]; ok {
			if intId, err := strconv.ParseInt(id.(string), 10, 64); err == nil {
				return payees, db.Find(&payees, intId).Error
			} else {
				return nil, err
			}
		}
		if name, ok := p.Args["name"]; ok {
			return payees, db.Find(&payees, "lower(name) = lower(?)", name.(string)).Error
		}
		return payees, db.Order("name").Find(&payees).Error
	},
}
