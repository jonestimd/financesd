package graphql

import (
	"strconv"

	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	"github.com/jonestimd/financesd/internal/dao/accountdao"
	"github.com/jonestimd/financesd/internal/model"
)

// Schema
var accountSchema = graphql.NewObject(graphql.ObjectConfig{
	Name:        "account",
	Description: "a financial account",
	Fields: graphql.Fields{
		"id":      &graphql.Field{Type: graphql.ID},
		"name":    &graphql.Field{Type: graphql.String},
		"company": &graphql.Field{Type: companySchema},
	},
})

var accountQueryFields = &graphql.Field{
	Type: graphql.NewList(accountSchema),
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
		if isSelected(accountQuery, p.Info, "company") {
			db = db.Preload("Company") // TODO make sure it only applies to account query and doesn't effect TX
		}
		if id, ok := p.Args["id"]; ok {
			if intId, err := strconv.ParseInt(id.(string), 10, 64); err == nil {
				account, err := accountdao.FindById(db, intId)
				return []*model.Account{account}, err
			} else {
				return nil, err
			}
		}
		if name, ok := p.Args["name"]; ok {
			account, err := accountdao.FindByName(db, name.(string))
			return []*model.Account{account}, err
		}
		return accountdao.GetAll(db)
	},
}
