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
			return getAccountByID(tx, id)
		}
		if nameArg, ok := p.Args["name"]; ok {
			name, _ := nameArg.(string)
			return getAccountsByName(tx, name)
		}
		return getAllAccounts(tx)
	},
}

type accountModel interface {
	GetCompany(tx *sql.Tx) (*model.Company, error)
}

func resolveCompany(p graphql.ResolveParams) (interface{}, error) {
	account := p.Source.(accountModel)
	tx := p.Context.Value(DbContextKey).(*sql.Tx)
	return account.GetCompany(tx)
}
