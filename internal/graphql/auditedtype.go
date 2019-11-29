package graphql

import (
	"github.com/graphql-go/graphql"
)

func addAudit(fields graphql.Fields) graphql.Fields {
	fields["version"] = &graphql.Field{Type: graphql.Int}
	fields["changeUser"] = &graphql.Field{Type: graphql.String, Resolve: nestedResolver("ChangeUser")}
	fields["changeDate"] = &graphql.Field{Type: graphql.DateTime, Resolve: nestedResolver("ChangeDate")}
	return fields
}
