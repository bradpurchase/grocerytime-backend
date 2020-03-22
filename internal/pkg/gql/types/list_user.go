package gql

import "github.com/graphql-go/graphql"

// ListUserType defines a graphql type for ListUser
var ListUserType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "ListUser",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"listId": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"userId": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"creator": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Boolean),
			},
			"active": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Boolean),
			},
			"createdAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"updatedAt": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	},
)
