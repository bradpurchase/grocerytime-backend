package gql

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/gql/resolvers"
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
			"itemsCount": &graphql.Field{
				Type:    graphql.Int,
				Resolve: resolvers.ListItemsCountResolver,
			},
			"trip": &graphql.Field{
				Type:    GroceryTripType,
				Resolve: resolvers.GroceryTripResolver,
			},
			"creator": &graphql.Field{
				Type:    BasicUserType,
				Resolve: resolvers.BasicUserResolver,
			},
			"listUsers": &graphql.Field{
				Type:    graphql.NewList(ListUserType),
				Resolve: resolvers.ListUsersResolver,
			},
		},
	},
)
