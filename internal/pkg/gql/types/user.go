package gql

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/gql/resolvers"
	"github.com/graphql-go/graphql"
)

// UserType defines a graphql type for User
var UserType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"email": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"firstName": &graphql.Field{
				Type: graphql.String,
			},
			"lastName": &graphql.Field{
				Type: graphql.String,
			},
			"createdAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"updatedAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"lists": &graphql.Field{
				Type:    graphql.NewList(ListType),
				Resolve: resolvers.ListsResolver,
			},
		},
	},
)
