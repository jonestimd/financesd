package schema

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MonsantoCo/mocka/v2"
	"github.com/graphql-go/graphql"
	"github.com/jonestimd/financesd/internal/model"
	"github.com/jonestimd/financesd/internal/sqltest"
	"github.com/stretchr/testify/assert"
)

func Test_securitySchema_Fields(t *testing.T) {
	symbol := "S1"
	now := time.Now()
	cost := float64(12.34)
	dividends := float64(23.45)
	security := &model.Security{
		Type:             "Stock",
		Shares:           96,
		FirstAcquired:    &now,
		CostBasis:        &cost,
		Dividends:        &dividends,
		TransactionCount: 69,
		Asset: model.Asset{
			ID:      1,
			Name:    "the asset",
			Type:    "Security",
			Scale:   6,
			Symbol:  &symbol,
			Version: 42,
			Audited: model.Audited{ChangeUser: "me", ChangeDate: &now},
		},
	}
	params := graphql.ResolveParams{Source: security}
	tests := []struct {
		field string
		value interface{}
	}{
		{"id", security.Asset.ID},
		{"assetType", security.Asset.Type},
		{"name", security.Asset.Name},
		{"scale", security.Asset.Scale},
		{"symbol", security.Asset.Symbol},
		{"version", security.Asset.Version},
		{"changeUser", security.Asset.Audited.ChangeUser},
		{"changeDate", security.Asset.Audited.ChangeDate},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s returns value", test.field), func(t *testing.T) {
			value, err := securitySchema.Fields()[test.field].Resolve(params)

			assert.Nil(t, err)
			assert.Equal(t, test.value, value)
		})
	}
}

func Test_securityQueryFields_Resolve_all(t *testing.T) {
	symbol := "S1"
	securities := []*model.Security{{AssetID: 42}}
	getAll := mocka.Function(t, &getAllSecurities, securities, nil)
	getByID := mocka.Function(t, &getSecurityByID, securities, nil)
	getBySymbol := mocka.Function(t, &getSecurityBySymbol, securities, nil)
	defer func() {
		getAll.Restore()
		getByID.Restore()
		getBySymbol.Restore()
	}()
	sqltest.TestInTx(t, func(mock sqlmock.Sqlmock, tx *sql.Tx) {
		tests := []struct {
			name     string
			argName  string
			argValue interface{}
			stub     *mocka.Stub
			stubArgs []interface{}
			err      bool
		}{
			{name: "returns all securities", stub: getAll, stubArgs: []interface{}{tx}},
			{name: "returns security with ID", argName: "id", argValue: 42, stub: getByID, stubArgs: []interface{}{tx, int64(42)}},
			{name: "returns security with symbol", argName: "symbol", argValue: symbol, stub: getBySymbol, stubArgs: []interface{}{tx, symbol}},
		}
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				params := newResolveParams(tx, securityQuery, newField("", "type"), newField("", "id"), newField("", "name")).
					addArg(test.argName, test.argValue)

				result, err := securityQueryFields.Resolve(params.ResolveParams)

				if test.err {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err)
					assert.Equal(t, securities, result)
					assert.Equal(t, 1, test.stub.CallCount())
					assert.Equal(t, test.stubArgs, test.stub.GetFirstCall().Arguments())
				}
			})
		}
	})
}
