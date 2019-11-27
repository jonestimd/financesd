package graphql

import (
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
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
	ParseValue: func(value interface{}) interface{} {
		result := &model.YesNo{}
		switch value := value.(type) {
		case bool:
			result.Set(value)
		case string:
			result.Set(strings.ToLower(value) == "true")
		case *string:
			result.Set(strings.ToLower(*value) == "true")
		default:
			result.Set(false)
		}
		return result
	},
	ParseLiteral: func(value ast.Value) interface{} {
		result := &model.YesNo{}
		switch value := value.(type) {
		case *ast.BooleanValue:
			result.Set(value.Value)
		case *ast.StringValue:
			result.Set(strings.ToLower(value.Value) == "true")
		default:
			result.Set(false)
		}
		return result
	},
})
