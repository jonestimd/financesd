package graphql

import (
	"github.com/graphql-go/graphql"
)

var yesNoType = graphql.NewScalar(graphql.ScalarConfig{
	Name: "boolean",
	Serialize: func(value interface{}) interface{} {
		switch value := value.(type) {
		case string:
			return value == "Y"
		case *string:
			return *value == "Y"
		}
		return nil
	},
	ParseValue:   graphql.Boolean.ParseValue,
	ParseLiteral: graphql.Boolean.ParseLiteral,
})
