package gql

import (
	"github.com/graphql-go/graphql"
)

// ListType defines a graphql type for List
var ListType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "List",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"userId": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"createdAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"updatedAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			// "listUsers": &graphql.Field{
			// 	Type: graphql.NewList(ListUserType),
			// 	Resolve: resolvers.ListUsersResolver,
			// },
		},
	},
)
