package schema

import "github.com/graphql-go/graphql"

// nonNullList creates a non-null list of non-null values.
func nonNullList(ofType graphql.Type) *graphql.NonNull {
	return graphql.NewNonNull(newList(ofType))
}

func newList(ofType graphql.Type) *graphql.List {
	return graphql.NewList(graphql.NewNonNull(ofType))
}

func asStrings(arg interface{}) []string {
	list := arg.([]interface{})
	values := make([]string, len(list))
	for i, value := range list {
		values[i] = value.(string)
	}
	return values
}

func asInts(arg interface{}) []int {
	list := arg.([]interface{})
	values := make([]int, len(list))
	for i, value := range list {
		values[i] = value.(int)
	}
	return values
}

var nonNullInt = graphql.NewNonNull(graphql.Int)
var nonNullFloat = graphql.NewNonNull(graphql.Float)
var nonNullString = graphql.NewNonNull(graphql.String)
var nonNullDate = graphql.NewNonNull(dateType)

var intList = graphql.NewList(nonNullInt)
var stringList = graphql.NewList(nonNullString)
