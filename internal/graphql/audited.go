package graphql

import (
	"reflect"

	"github.com/graphql-go/graphql"
)

func resolveChangeUser(p graphql.ResolveParams) (interface{}, error) {
	value := reflect.ValueOf(p.Source)
	field := value.Elem().FieldByName("ChangeUser")
	return field.String(), nil
}

func resolveChangeDate(p graphql.ResolveParams) (interface{}, error) {
	value := reflect.ValueOf(p.Source)
	field := value.Elem().FieldByName("ChangeDate")
	return field.Interface(), nil
}

func addAudit(fields graphql.Fields) graphql.Fields {
	fields["version"] = &graphql.Field{Type: graphql.Int}
	fields["changeUser"] = &graphql.Field{Type: graphql.String, Resolve: resolveChangeUser}
	fields["changeDate"] = &graphql.Field{Type: graphql.DateTime, Resolve: resolveChangeDate}
	return fields
}
