package schema

import (
	"reflect"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

type selectionPredicate func(field *ast.Field) bool

func findQuery(info graphql.ResolveInfo) *ast.Field {
	return findField(info.Operation.GetSelectionSet().Selections, byAlias(info.Path.Key))
}

func byAlias(alias interface{}) selectionPredicate {
	return func(field *ast.Field) bool {
		if field.Alias != nil {
			return field.Alias.Value == alias
		}
		return field.Name.Value == alias
	}
}

func byName(name string) selectionPredicate {
	return func(field *ast.Field) bool {
		return field.Name.Value == name
	}
}

func findField(selections []ast.Selection, predicate selectionPredicate) *ast.Field {
	for _, field := range selections {
		switch field := field.(type) {
		case *ast.Field:
			if predicate(field) {
				return field
			}
		}
	}
	return nil
}

// nestedResolver returns a resolver for the nested field of the source specified by the chain of fieldNames
func nestedResolver(fieldNames ...string) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		value := reflect.ValueOf(p.Source)
		for _, fieldName := range fieldNames {
			if value.Kind() == reflect.Interface || value.Kind() == reflect.Ptr {
				value = value.Elem()
			}
			value = value.FieldByName(fieldName)
		}
		return value.Interface(), nil
	}
}
