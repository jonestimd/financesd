package graphql

import (
	"reflect"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

type selectionPredicate func(field *ast.Field) bool

// isPathSelected checks if a property path has been selected for a query.
func isPathSelected(info graphql.ResolveInfo, path ...string) bool {
	query := findField(info.Operation.GetSelectionSet().Selections, byAlias(info.Path.Key))
	fields := query.GetSelectionSet().Selections
	for _, name := range path {
		field := findField(fields, byName(name))
		if field == nil {
			return false
		}
		fields = field.GetSelectionSet().Selections
	}
	return true
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

// nestedResolver returns a resolver for a nested field
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
