package schema

import (
	"database/sql"
	"strconv"

	"github.com/graphql-go/graphql"
	"github.com/jonestimd/financesd/internal/model"
)

// Schema
var companySchema = graphql.NewObject(graphql.ObjectConfig{
	Name:        "company",
	Description: "a financial company",
	Fields: addAudit(graphql.Fields{
		"id":   &graphql.Field{Type: graphql.ID},
		"name": &graphql.Field{Type: graphql.String},
	}),
})

func companyQueryFields() *graphql.Field {
	companySchema.AddFieldConfig("accounts", &graphql.Field{Type: graphql.NewList(accountSchema), Resolve: resolveAccounts})
	return &graphql.Field{
		Type: graphql.NewList(companySchema),
		Args: graphql.FieldConfigArgument{
			"id":   {Type: graphql.ID, Description: "company ID"},
			"name": {Type: graphql.String, Description: "unique company name"},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			tx := p.Context.Value(DbContextKey).(*sql.Tx)
			if idArg, ok := p.Args["id"]; ok {
				id, err := strconv.ParseInt(idArg.(string), 10, 64)
				if err != nil {
					return nil, err
				}
				return getCompanyByID(tx, id)
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
	GetAccounts(tx *sql.Tx) ([]*model.Account, error)
}

func resolveAccounts(p graphql.ResolveParams) (interface{}, error) {
	company := p.Source.(companyModel)
	tx := p.Context.Value(DbContextKey).(*sql.Tx)
	return company.GetAccounts(tx)
}

var addCompaniesFields = &graphql.Field{
	Type: newList(companySchema),
	Args: graphql.FieldConfigArgument{
		"names": {Type: newList(graphql.String), Description: "unique company names"},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		tx := p.Context.Value(DbContextKey).(*sql.Tx)
		user := p.Context.Value(UserKey).(string)
		names := asStrings(p.Args["names"])
		return addCompanies(tx, names, user)
	},
}
