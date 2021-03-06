package schema

import (
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

const dateFormat = "2006-01-02"

var dateType = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "Date",
	Description: "Date contains a calendar date (without time of day)",
	Serialize: func(value interface{}) interface{} {
		switch value := value.(type) {
		case time.Time:
			return value.Format(dateFormat)
		case *time.Time:
			return value.Format(dateFormat)
		}
		return nil
	},
	ParseValue: func(value interface{}) interface{} {
		switch value := value.(type) {
		case string:
			if result, err := time.Parse(dateFormat, value); err == nil {
				return result
			}
		case *string:
			if result, err := time.Parse(dateFormat, *value); err == nil {
				return result
			}
		}
		return nil
	},
	ParseLiteral: func(value ast.Value) interface{} {
		switch value := value.(type) {
		case *ast.StringValue:
			if result, err := time.Parse(dateFormat, value.Value); err == nil {
				return result
			}
		}
		return nil
	},
})
