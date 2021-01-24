package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/jonestimd/financesd/internal/model"
)

const accountsRootKey = "accountsByID"
const companyIDsRootKey = "companyIDs"
const companiesRootKey = "companiesByID"
const accountTxKey = "accountTxKey"

type companyIDs interface {
	Result() interface{}
	Error() error
	CompanyIDs() []int64
}

func rootValue(info graphql.ResolveInfo) map[string]interface{} {
	return info.RootValue.(map[string]interface{})
}

func addCompanyIDsToRoot(info graphql.ResolveInfo, result companyIDs) (interface{}, error) {
	if result.Error() != nil {
		return nil, result.Error()
	}
	root := rootValue(info)
	root[companyIDsRootKey] = result.CompanyIDs()
	return result.Result(), nil
}

func addAccountsToRoot(info graphql.ResolveInfo, result *model.Accounts) ([]*model.Account, error) {
	if result.Error() != nil {
		return nil, result.Error()
	}
	root := rootValue(info)
	root[companyIDsRootKey] = result.CompanyIDs()
	root[accountsRootKey] = result.ByID()
	return result.Accounts, nil
}
