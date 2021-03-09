package schema

import "github.com/graphql-go/graphql"

// newList creates a non-null list of non-null values.
func newList(ofType graphql.Type) *graphql.NonNull {
	return graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(ofType)))
}

func asStrings(arg interface{}) []string {
	list := arg.([]interface{})
	values := make([]string, len(list))
	for i, value := range list {
		values[i] = value.(string)
	}
	return values
}
