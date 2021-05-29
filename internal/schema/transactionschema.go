package schema

import (
	"database/sql"
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/jonestimd/financesd/internal/database"
	"github.com/jonestimd/financesd/internal/domain"
)

func getDetailFields() graphql.Fields {
	return graphql.Fields{
		"id":                    &graphql.Field{Type: nonNullInt},
		"transactionCategoryId": &graphql.Field{Type: graphql.Int},
		"transactionGroupId":    &graphql.Field{Type: graphql.Int},
		"memo":                  &graphql.Field{Type: graphql.String},
		"amount":                &graphql.Field{Type: nonNullFloat},
		"assetQuantity":         &graphql.Field{Type: graphql.Float},
		"exchangeAssetId":       &graphql.Field{Type: graphql.Int},
	}
}

func getDetailSchema(name string, relatedField string, fieldType graphql.Output, resolve graphql.FieldResolveFn) *graphql.Object {
	detailFields := getDetailFields()
	detailFields[relatedField] = &graphql.Field{Type: fieldType, Resolve: resolve}
	return graphql.NewObject(graphql.ObjectConfig{
		Name:        name,
		Description: "a detail of a financial transaction",
		Fields:      addAudit(detailFields),
	})
}

func getTxFields() graphql.Fields {
	return graphql.Fields{
		"id":              &graphql.Field{Type: nonNullInt},
		"date":            &graphql.Field{Type: nonNullDate},
		"memo":            &graphql.Field{Type: graphql.String},
		"referenceNumber": &graphql.Field{Type: graphql.String},
		"cleared":         &graphql.Field{Type: yesNoType},
		"accountId":       &graphql.Field{Type: nonNullInt},
		"payeeId":         &graphql.Field{Type: graphql.Int},
		"securityId":      &graphql.Field{Type: graphql.Int},
	}
}

func getTxSchemaConfig(name string) graphql.ObjectConfig {
	return graphql.ObjectConfig{
		Description: "a financial transaction",
		Name:        name,
		Fields:      addAudit(getTxFields()),
	}
}

func getTxSchema() *graphql.Object {
	relatedTxSchema := graphql.NewObject(getTxSchemaConfig("relatedTransaction"))
	relatedDetailSchema := getDetailSchema("relatedTransactionDetail", "transaction", relatedTxSchema, resolveRelatedTransaction)
	detailSchema := getDetailSchema("transactionDetail", "relatedDetail", relatedDetailSchema, resolveRelatedDetail)
	txSchema := graphql.NewObject(getTxSchemaConfig("transaction"))
	txSchema.AddFieldConfig("details", &graphql.Field{Type: graphql.NewList(detailSchema), Resolve: resolveDetails})
	return txSchema
}

var txList = newList(getTxSchema())

var transactionQueryFields = &graphql.Field{
	Type: txList,
	Args: graphql.FieldConfigArgument{
		"accountId": {Type: graphql.NewNonNull(graphql.Int), Description: "account ID"},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		tx := p.Context.Value(DbContextKey).(*sql.Tx)
		accountID := p.Args["accountId"].(int)
		return getAccountTransactions(tx, int64(accountID)), nil
	},
}

type txModel interface {
	GetDetails(tx *sql.Tx) []*domain.TransactionDetail
}

var _ txModel = (*domain.Transaction)(nil)

func resolveDetails(p graphql.ResolveParams) (interface{}, error) {
	if transaction, ok := p.Source.(txModel); ok {
		tx := p.Context.Value(DbContextKey).(*sql.Tx)
		return transaction.GetDetails(tx), nil
	}
	return nil, errors.New("invalid source")
}

type detailModel interface {
	GetRelatedDetail(tx *sql.Tx) *domain.TransactionDetail
	GetRelatedTransaction(tx *sql.Tx) *domain.Transaction
}

var _ detailModel = (*domain.TransactionDetail)(nil)

func resolveRelatedDetail(p graphql.ResolveParams) (interface{}, error) {
	if detail, ok := p.Source.(detailModel); ok {
		tx := p.Context.Value(DbContextKey).(*sql.Tx)
		return detail.GetRelatedDetail(tx), nil
	}
	return nil, errors.New("invalid source")
}

func resolveRelatedTransaction(p graphql.ResolveParams) (interface{}, error) {
	if detail, ok := p.Source.(detailModel); ok {
		tx := p.Context.Value(DbContextKey).(*sql.Tx)
		return detail.GetRelatedTransaction(tx), nil
	}
	return nil, errors.New("invalid source")
}

func getDetailInput(action string) *graphql.InputObjectFieldConfig {
	fields := graphql.InputObjectConfigFieldMap{
		"transferAccountId": &graphql.InputObjectFieldConfig{Type: graphql.Int},
		"version":           &graphql.InputObjectFieldConfig{Type: graphql.Int},
	}
	for name, field := range getDetailFields() {
		fieldType := field.Type
		if action == "update" && (name == "id" || name == "amount") {
			fieldType = field.Type.(*graphql.NonNull).OfType
		}
		fields[name] = &graphql.InputObjectFieldConfig{Type: fieldType}
	}
	return &graphql.InputObjectFieldConfig{
		Description: "Update data for transaction details.",
		Type: graphql.NewList(graphql.NewNonNull(graphql.NewInputObject(graphql.InputObjectConfig{
			Name: action + "TransactionDetailInput",
			Description: "To update, provide **id**, **version** and fields to modify." +
				" To add, provide at least **amount** (without **id** or **version**)." +
				" To delete, provide just **id** and **version**.",
			Fields: fields,
		}))),
	}
}

func getTxInput(action string) *graphql.InputObject {
	fields := graphql.InputObjectConfigFieldMap{"details": getDetailInput(action)}
	for name, field := range getTxFields() {
		if (action != "add" || name != "id") && name != "accountId" {
			fields[name] = &graphql.InputObjectFieldConfig{Type: field.Type}
		}
	}
	if action == "update" {
		fields["date"].Type = dateType
		fields["accountId"] = &graphql.InputObjectFieldConfig{Type: graphql.Int}
		fields["version"] = &graphql.InputObjectFieldConfig{Type: nonNullInt, Description: "The current version of the transaction."}
	}
	return graphql.NewInputObject(graphql.InputObjectConfig{
		Name:   action + "TransactionInput",
		Fields: fields,
	})
}

var updateTxFields = &graphql.Field{
	Type:        txList,
	Description: "Add, update and/or delete transactions.",
	Args: graphql.FieldConfigArgument{
		"accountId": {Type: nonNullInt, Description: "ID of account for transactions."},
		"add":       {Type: newList(getTxInput("add")), Description: "Transactions to add."},
		"update":    {Type: newList(getTxInput("update")), Description: "Changes to be made to existing transactions."},
		"delete":    {Type: idVersionList, Description: "IDs of transactions to delete."},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		transactions := []*domain.Transaction{}
		tx := p.Context.Value(DbContextKey).(*sql.Tx)
		user := p.Context.Value(UserKey).(string)
		if ids, ok := p.Args["delete"]; ok {
			deleteTransactions(tx, asMaps(ids))
		}
		ids := make([]int64, 0)
		if updates, ok := p.Args["update"]; ok {
			ids = updateTransactions(tx, asMaps(updates, "details"), user)
		}
		if inserts, ok := p.Args["add"]; ok {
			accountID := database.InputObject(p.Args).RequireInt("accountId")
			ids = append(ids, insertTransactions(tx, accountID, asMaps(inserts, "details"), user)...)
		}
		if len(ids) > 0 {
			transactions = getTransactionsByIDs(tx, ids)
		}
		return transactions, nil
	},
}
