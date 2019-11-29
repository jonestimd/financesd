package graphql

import (
	"strconv"

	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	"github.com/jonestimd/financesd/internal/model"
)

// Schema
var companySchema = graphql.NewObject(graphql.ObjectConfig{
	Name:        "company",
	Description: "a financial company",
	Fields: addAudit(graphql.Fields{
		"id":   &graphql.Field{Type: graphql.ID},
		"name": &graphql.Field{Type: graphql.String},
		// add later in Schema(): "accounts": &graphql.Field{Type: graphql.NewList(accountSchema())},
	}),
})

var companyQueryFields = &graphql.Field{
	Type: graphql.NewList(companySchema),
	Args: map[string]*graphql.ArgumentConfig{
		"id":   &graphql.ArgumentConfig{Type: graphql.ID, Description: "company ID"},
		"name": &graphql.ArgumentConfig{Type: graphql.String, Description: "unique company name"},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		db := p.Context.Value(DbContextKey).(*gorm.DB)
		companies := make([]*model.Company, 0)
		if isPathSelected(p.Info, "accounts") {
			db = db.Preload("Accounts") // TODO make sure it only applies to company query and doesn't effect TX
		}
		if id, ok := p.Args["id"]; ok {
			if intId, err := strconv.ParseInt(id.(string), 10, 64); err == nil {
				return companies, db.Find(&companies, intId).Error
			} else {
				return nil, err
			}
		}
		if name, ok := p.Args["name"]; ok {
			return companies, db.Find(&companies, "lower(name) = lower(?)", name.(string)).Error
		}
		return companies, db.Order("name").Find(&companies).Error
	},
}
