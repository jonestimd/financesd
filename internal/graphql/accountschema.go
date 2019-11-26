package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	"github.com/jonestimd/financesd/internal/dao/accountdao"
	"github.com/jonestimd/financesd/internal/dao/companydao"
	"github.com/jonestimd/financesd/internal/model"
)

// Schema
var accountFields = graphql.NewObject(graphql.ObjectConfig{
	Name:        "account",
	Description: "a financial account",
	Fields: graphql.Fields{
		"id":   &graphql.Field{Type: graphql.ID},
		"name": &graphql.Field{Type: graphql.String},
		"company": &graphql.Field{
			Type: companyFields,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				account := p.Source.(*model.Account)
				if account.CompanyID.Valid {
					db := p.Context.Value(DbContextKey).(*gorm.DB)
					return companydao.FindById(db, account.CompanyID.Int64)
				}
				return nil, nil
			},
		},
	},
})

var accountQueryFields = &graphql.Field{
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
		db := p.Context.Value(DbContextKey).(*gorm.DB)
		if id, ok := p.Args["id"]; ok {
			account, err := accountdao.FindById(db, id.(int64))
			return []*model.Account{account}, err
		}
		if name, ok := p.Args["name"]; ok {
			account, err := accountdao.FindByName(db, name.(string))
			return []*model.Account{account}, err
		}
		return accountdao.GetAll(db)
	},
}
