package schema

import (
	"github.com/graphql-go/graphql"
)

func addAudit(fields graphql.Fields) graphql.Fields {
	fields["version"] = &graphql.Field{Type: graphql.Int}
	fields["changeUser"] = &graphql.Field{Type: graphql.String, Resolve: nestedResolver("Audited", "ChangeUser")}
	fields["changeDate"] = &graphql.Field{Type: graphql.String, Resolve: nestedResolver("Audited", "ChangeDate")}
	return fields
}
