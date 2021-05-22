package schema

import (
	"database/sql"

	"github.com/graphql-go/graphql"
	"github.com/jonestimd/financesd/internal/domain"
)

// Schema
var companySchema = graphql.NewObject(graphql.ObjectConfig{
	Name:        "company",
	Description: "a financial company",
	Fields: addAudit(graphql.Fields{
		"id":   &graphql.Field{Type: graphql.Int},
		"name": &graphql.Field{Type: graphql.String},
	}),
})

var companyList = nonNullList(companySchema)

func companyQueryFields() *graphql.Field {
	companySchema.AddFieldConfig("accounts", &graphql.Field{Type: graphql.NewList(accountSchema), Resolve: resolveAccounts})
	return &graphql.Field{
		Type: companyList,
		Args: graphql.FieldConfigArgument{
			"id":   {Type: graphql.Int, Description: "company ID"},
			"name": {Type: graphql.String, Description: "unique company name"},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			tx := p.Context.Value(DbContextKey).(*sql.Tx)
			if idArg, ok := p.Args["id"]; ok {
				id := idArg.(int)
				return getCompanyByID(tx, int64(id))
			}
			if nameArg, ok := p.Args["name"]; ok {
				name, _ := nameArg.(string)
				return getCompanyByName(tx, name)
			}
			return getAllCompanies(tx)
		},
	}
}

type companyModel interface {
	GetAccounts(tx *sql.Tx) ([]*domain.Account, error)
}

func resolveAccounts(p graphql.ResolveParams) (interface{}, error) {
	company := p.Source.(companyModel)
	tx := p.Context.Value(DbContextKey).(*sql.Tx)
	return company.GetAccounts(tx)
}

var companyInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "companyInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"id":      {Type: nonNullInt, Description: "ID of the company to update."},
		"name":    {Type: nonNullString, Description: "New name for the company."},
		"version": {Type: nonNullInt, Description: "Current version of the company."},
	},
})

var updateCompaniesFields = &graphql.Field{
	Type:        companyList,
	Description: "Add, update and/or delete companies.",
	Args: graphql.FieldConfigArgument{
		"add":    {Type: stringList, Description: "Unique names of companies to add."},
		"update": {Type: newList(companyInput), Description: "Changes to be made to existing companies."},
		"delete": {Type: intList, Description: "IDs of companies to delete."},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		var err error
		companies := make([]*domain.Company, 0)
		tx := p.Context.Value(DbContextKey).(*sql.Tx)
		user := p.Context.Value(UserKey).(string)
		if ids, ok := p.Args["delete"]; ok {
			if _, err = deleteCompanies(tx, asInts(ids)); err != nil {
				return nil, err
			}
		}
		if updates, ok := p.Args["update"]; ok {
			if companies, err = updateCompanies(tx, updates, user); err != nil {
				return nil, err
			}
		}
		if names, ok := p.Args["add"]; ok {
			if added, err := addCompanies(tx, asStrings(names), user); err != nil {
				return nil, err
			} else {
				companies = append(companies, added...)
			}
		}
		return companies, nil
	},
}
