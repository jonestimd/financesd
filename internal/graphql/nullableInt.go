package graphql

import (
	"database/sql"
	"strconv"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

var nullableInt = graphql.NewScalar(graphql.ScalarConfig{
	Name: "int",
	Description: "The `int` scalar type represents non-fractional signed whole numeric " +
		"values. Int can represent values between -(2^31) and 2^31 - 1. ",
	Serialize: func(value interface{}) interface{} {
		switch value := value.(type) {
		case sql.NullInt64:
			if value.Valid {
				return value.Int64
			}
		case *sql.NullInt64:
			if value.Valid {
				return value.Int64
			}
		}
		return nil
	},
	ParseValue: func(value interface{}) interface{} {
		switch value := value.(type) {
		case int:
			return &sql.NullInt64{Int64: int64(value), Valid: true}
		case int64:
			return &sql.NullInt64{Int64: value, Valid: true}
		case string:
			if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
				return &sql.NullInt64{Int64: intValue, Valid: true}
			}
		case *string:
			if intValue, err := strconv.ParseInt(*value, 10, 64); err == nil {
				return &sql.NullInt64{Int64: intValue, Valid: true}
			}
		}
		return &sql.NullInt64{Valid: false}
	},
	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch valueAST := valueAST.(type) {
		case *ast.IntValue:
			if intValue, err := strconv.ParseInt(valueAST.Value, 10, 64); err == nil {
				return &sql.NullInt64{Int64: intValue, Valid: true}
			}
		}
		return &sql.NullInt64{Valid: false}
	},
})
