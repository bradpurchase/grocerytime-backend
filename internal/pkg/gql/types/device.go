package gql

import "github.com/graphql-go/graphql"

// DeviceType defines a graphql type for Device
var DeviceType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Device",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"userId": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"token": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"createdAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"updatedAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"deletedAt": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	},
)
