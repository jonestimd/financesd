package graphql

import (
	"database/sql"
	"strconv"

	"github.com/graphql-go/graphql"
	"github.com/jonestimd/financesd/internal/model"
)

// Schema
var accountSchema = graphql.NewObject(graphql.ObjectConfig{
	Name:        "account",
	Description: "a financial account",
	Fields: addAudit(graphql.Fields{
		"id":               &graphql.Field{Type: graphql.ID},
		"companyId":        &graphql.Field{Type: graphql.Int},
		"company":          &graphql.Field{Type: companySchema, Resolve: resolveCompany},
		"name":             &graphql.Field{Type: graphql.String},
		"description":      &graphql.Field{Type: graphql.String},
		"accountNo":        &graphql.Field{Type: graphql.String},
		"type":             &graphql.Field{Type: graphql.String}, // TODO enum?
		"closed":           &graphql.Field{Type: yesNoType},
		"currencyId":       &graphql.Field{Type: graphql.Int},
		"transactionCount": &graphql.Field{Type: graphql.Int},
		"balance":          &graphql.Field{Type: graphql.String},
	}),
})

var accountQueryFields = &graphql.Field{
	Type: graphql.NewList(accountSchema),
	Args: map[string]*graphql.ArgumentConfig{
		"id":   {Type: graphql.ID, Description: "account ID"},
		"name": {Type: graphql.String, Description: "unique account name"},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		tx := p.Context.Value(DbContextKey).(*sql.Tx)
		if idArg, ok := p.Args["id"]; ok {
			id, err := strconv.ParseInt(idArg.(string), 10, 64)
			if err != nil {
				return nil, err
			}
			return addCompanyIDsToRoot(p.Info, getAccountByID(tx, id))
		}
		if nameArg, ok := p.Args["name"]; ok {
			name, _ := nameArg.(string)
			return addCompanyIDsToRoot(p.Info, getAccountsByName(tx, name))
		}
		return addAccountsToRoot(p.Info, getAllAccounts(tx))
	},
}

func resolveCompany(p graphql.ResolveParams) (interface{}, error) {
	account := p.Source.(*model.Account)
	if account.CompanyID == nil {
		return nil, nil
	}
	tx := p.Context.Value(DbContextKey).(*sql.Tx)
	rootValue := p.Info.RootValue.(map[string]interface{})
	if _, ok := rootValue[companiesRootKey]; !ok {
		ids := rootValue[companyIDsRootKey].([]int64)
		byID, err := getCompaniesByIDs(tx, ids)
		if err != nil {
			return nil, err
		}
		rootValue[companiesRootKey] = byID
	}
	companiesByID := rootValue[companiesRootKey].(map[int64]*model.Company)
	return companiesByID[*account.CompanyID], nil
}
