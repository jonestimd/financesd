package model

import "github.com/graphql-go/graphql"

func replaceSource(p graphql.ResolveParams, source interface{}) graphql.ResolveParams {
	return graphql.ResolveParams{
		Source:  source,
		Args:    p.Args,
		Info:    p.Info,
		Context: p.Context,
	}
}
