package gql

import (
	"github.com/graphql-go/graphql"
)

// StoreStapleItemType defines a graphql type for Item
var StoreStapleItemType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "StoreStapleItem",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"storeId": &graphql.Field{
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
		},
	},
)
