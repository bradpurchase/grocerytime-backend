package gql

import "github.com/graphql-go/graphql"

var StoreCategoryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "StoreCategory",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"storeId": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"categoryId": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
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
