package graphql

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MonsantoCo/mocka/v2"
	"github.com/jonestimd/financesd/internal/model"
	"github.com/jonestimd/financesd/internal/sqltest"
	"github.com/stretchr/testify/assert"
)

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
			{name: "returns security with ID", argName: "id", argValue: "42", stub: getByID, stubArgs: []interface{}{tx, int64(42)}},
			{name: "returns security with symbol", argName: "symbol", argValue: symbol, stub: getBySymbol, stubArgs: []interface{}{tx, symbol}},
			{name: "returns error for invalid ID", argName: "id", argValue: "abc", stub: getByID, err: true},
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
