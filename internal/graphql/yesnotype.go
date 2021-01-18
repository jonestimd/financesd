package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/jonestimd/financesd/internal/model"
)

var yesNoType = graphql.NewScalar(graphql.ScalarConfig{
	Name: "boolean",
	Serialize: func(value interface{}) interface{} {
		switch value := value.(type) {
		case model.YesNo:
			return value.Get()
		case *model.YesNo:
			return value.Get()
		}
		return nil
	},
	ParseValue:   graphql.Boolean.ParseValue,
	ParseLiteral: graphql.Boolean.ParseLiteral,
})
