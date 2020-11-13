package gql

import "github.com/graphql-go/graphql"

// StoreUserPreferenceType defines a graphql type for StoreUserPreference
var StoreUserPreferenceType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "StoreUserPreference",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"storeUserId": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"defaultStore": &graphql.Field{
				Type: graphql.Boolean,
			},
			"notifications": &graphql.Field{
				Type: graphql.Boolean,
			},
			"createdAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"updatedAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"deleteAt": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	},
)
