package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	"github.com/jonestimd/financesd/internal/dao/accountdao"
)

// Schema
var accountFields = graphql.NewObject(graphql.ObjectConfig{
	Name:        "account",
	Description: "a financial account",
	Fields: graphql.Fields{
		"id":   &graphql.Field{Type: graphql.ID},
		"name": &graphql.Field{Type: graphql.String},
	},
})

var accountList = graphql.Fields{
	"accounts": &graphql.Field{
		Type: graphql.NewList(accountFields),
		Args: map[string]*graphql.ArgumentConfig{
			"id": &graphql.ArgumentConfig{
				Type:        graphql.ID,
				Description: "account ID",
			},
			"name": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "unique account name",
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			db := p.Source.(map[string]interface{})["db"].(*gorm.DB)
			if id, ok := p.Args["id"]; ok {
				return accountdao.FindById(db, id.(int))
			}
			if name, ok := p.Args["name"]; ok {
				return accountdao.FindByName(db, name.(string))
			}
			return accountdao.GetAll(db)
		},
	},
}
