package graphql

import (
	"database/sql"
	"errors"
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
		// add in Schema(): "accounts": &graphql.Field{Type: graphql.NewList(accountSchema()), Resolve: resolveAccounts},
	}),
})

var companyQueryFields = &graphql.Field{
	Type: graphql.NewList(companySchema),
	Args: map[string]*graphql.ArgumentConfig{
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
			return addCompanyIDsToRoot(p.Info, getCompanyByID(tx, id))
		}
		if nameArg, ok := p.Args["name"]; ok {
			name, _ := nameArg.(string)
			return addCompanyIDsToRoot(p.Info, getCompanyByName(tx, name))
		}
		return addCompanyIDsToRoot(p.Info, getAllCompanies(tx))
	},
}

func resolveAccounts(p graphql.ResolveParams) (interface{}, error) {
	company := p.Source.(*model.Company)
	rootValue := p.Info.RootValue.(map[string]interface{})
	if _, ok := rootValue[accountsRootKey]; !ok {
		if ids, ok := rootValue[companyIDsRootKey]; ok {
			tx := p.Context.Value(DbContextKey).(*sql.Tx)
			if _, err := addAccountsToRoot(p.Info, getAccountsByCompanyIDs(tx, ids.([]int64))); err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("no company IDs in context")
		}
	}
	accountsByID := rootValue[accountsRootKey].(map[int64]*model.Account)
	accounts := make([]*model.Account, 0)
	for _, account := range accountsByID {
		if account.CompanyID != nil && *account.CompanyID == company.ID {
			accounts = append(accounts, account)
		}
	}
	return accounts, nil
}
