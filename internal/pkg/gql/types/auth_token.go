package gql

import (
	"github.com/graphql-go/graphql"
)

// AuthTokenType defines a graphql type for AuthToken
var AuthTokenType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "AuthToken",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"clientId": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"userId": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"accessToken": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"refreshToken": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"expiresIn": &graphql.Field{
				Type: graphql.DateTime,
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
