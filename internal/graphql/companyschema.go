package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	"github.com/jonestimd/financesd/internal/dao/companydao"
	"github.com/jonestimd/financesd/internal/model"
)

// Schema
var companyFields = graphql.NewObject(graphql.ObjectConfig{
	Name:        "company",
	Description: "a financial company",
	Fields: graphql.Fields{
		"id":   &graphql.Field{Type: graphql.ID},
		"name": &graphql.Field{Type: graphql.String},
	},
})

var companyQueryFields = &graphql.Field{
	Type: graphql.NewList(companyFields),
	Args: map[string]*graphql.ArgumentConfig{
		"id": &graphql.ArgumentConfig{
			Type:        graphql.ID,
			Description: "company ID",
		},
		"name": &graphql.ArgumentConfig{
			Type:        graphql.String,
			Description: "unique company name",
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		db := p.Context.Value(DbContextKey).(*gorm.DB)
		if id, ok := p.Args["id"]; ok {
			company, err := companydao.FindById(db, id.(int64))
			return []*model.Company{company}, err
		}
		if name, ok := p.Args["name"]; ok {
			company, err := companydao.FindByName(db, name.(string))
			return []*model.Company{company}, err
		}
		return companydao.GetAll(db)
	},
}
