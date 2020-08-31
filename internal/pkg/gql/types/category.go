package gql

import "github.com/graphql-go/graphql"

var CategoryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Category",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
	},
)
