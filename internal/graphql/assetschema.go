package graphql

import (
	"github.com/graphql-go/graphql"
)

func addAsset(fields graphql.Fields) graphql.Fields {
	fields["id"] = &graphql.Field{Type: graphql.ID, Resolve: nestedResolver("Asset", "ID")}
	fields["name"] = &graphql.Field{Type: graphql.String, Resolve: nestedResolver("Asset", "Name")}
	fields["scale"] = &graphql.Field{Type: graphql.Int, Resolve: nestedResolver("Asset", "Scale")}
	fields["symbol"] = &graphql.Field{Type: graphql.String, Resolve: nestedResolver("Asset", "Symbol")}
	fields["version"] = &graphql.Field{Type: graphql.String, Resolve: nestedResolver("Asset", "Version")}
	fields["changeUser"] = &graphql.Field{Type: graphql.String, Resolve: nestedResolver("Asset", "ChangeUser")}
	fields["changeDate"] = &graphql.Field{Type: graphql.String, Resolve: nestedResolver("Asset", "ChangeDate")}
	return fields
}
