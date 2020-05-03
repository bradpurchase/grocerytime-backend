package gql

import (
	"github.com/graphql-go/graphql"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/gql/resolvers"
	gql "github.com/bradpurchase/grocerytime-backend/internal/pkg/gql/types"
)

// Schema defines a graphql.Schema instance
var Schema graphql.Schema

func init() {
	// Define the root query type
	queryType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "RootQuery",
			Fields: graphql.Fields{
				"me": &graphql.Field{
					Type:        gql.UserType,
					Description: "Retrieve the current user",
					Resolve:     resolvers.AuthenticatedUserResolver,
				},
				"list": &graphql.Field{
					Type:        gql.ListType,
					Description: "Retrieve a list and its items",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
					},
					Resolve: resolvers.ListResolver,
				},
			},
		},
	)

	// Define the root mutation type
	mutationType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Mutation",
			Fields: graphql.Fields{
				"signup": &graphql.Field{
					Type:        gql.AuthTokenType,
					Description: "Create a new user account",
					Args: graphql.FieldConfigArgument{
						"email": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"password": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
					},
					Resolve: resolvers.SignupResolver,
				},
				"token": &graphql.Field{
					Type:        gql.AuthTokenType,
					Description: "Retrieve an access token",
					Args: graphql.FieldConfigArgument{
						"grantType": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"email": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"password": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"refreshToken": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: resolvers.TokenResolver,
				},
				"createList": &graphql.Field{
					Type:        gql.ListType,
					Description: "Create a list",
					Args: graphql.FieldConfigArgument{
						"name": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
					},
					Resolve: resolvers.CreateListResolver,
				},
				"deleteList": &graphql.Field{
					Type:        gql.ListType,
					Description: "Delete a list",
					Args: graphql.FieldConfigArgument{
						"listId": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
					},
					Resolve: resolvers.DeleteListResolver,
				},
				"updateList": &graphql.Field{
					Type:        gql.ListType,
					Description: "Update a list",
					Args: graphql.FieldConfigArgument{
						"listId": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
						"name": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: resolvers.UpdateListResolver,
				},
				"addListUser": &graphql.Field{
					Type:        gql.ListUserType,
					Description: "Add a user to a list",
					Args: graphql.FieldConfigArgument{
						"listId": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
						"email": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
					},
					Resolve: resolvers.AddListUserResolver,
				},
				"addItemToList": &graphql.Field{
					Type:        gql.ItemType,
					Description: "Add an item to a list",
					Args: graphql.FieldConfigArgument{
						"listId": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
						"name": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"quantity": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
					},
					Resolve: resolvers.AddItemResolver,
				},
				"deleteItem": &graphql.Field{
					Type:        gql.ItemType,
					Description: "Remove an item from a list",
					Args: graphql.FieldConfigArgument{
						"itemId": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
					},
					Resolve: resolvers.DeleteItemResolver,
				},
				"updateItem": &graphql.Field{
					Type:        gql.ItemType,
					Description: "Updates the properties of an item in a list",
					Args: graphql.FieldConfigArgument{
						"itemId": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
						"completed": &graphql.ArgumentConfig{
							Type: graphql.Boolean,
						},
					},
					Resolve: resolvers.UpdateItemResolver,
				},
			},
		},
	)

	var err error
	Schema, err = graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    queryType,
			Mutation: mutationType,
		},
	)
	if err != nil {
		panic(err)
	}
}
