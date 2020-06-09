package gql

import (
	"github.com/graphql-go/graphql"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/gql/resolvers"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/gql/subscriptions"
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
				"sharableList": &graphql.Field{
					Type:        gql.ListType,
					Description: "Retrieve basic info about a list for sharing",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
					},
					Resolve: resolvers.SharableListResolver,
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
				"joinList": &graphql.Field{
					Type:        gql.ListUserType,
					Description: "Join a list via share link",
					Args: graphql.FieldConfigArgument{
						"listId": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
					},
					Resolve: resolvers.JoinListResolver,
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
						"name": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"completed": &graphql.ArgumentConfig{
							Type: graphql.Boolean,
						},
						"quantity": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
						"position": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
					},
					Resolve: resolvers.UpdateItemResolver,
				},
				"reorderItem": &graphql.Field{
					Type:        gql.ListType,
					Description: "Updates the order of an item in a list",
					Args: graphql.FieldConfigArgument{
						"itemId": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
						"position": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.Int),
						},
					},
					Resolve: resolvers.ReorderItemResolver,
				},
				"addItemToTrip": &graphql.Field{
					Type:        gql.ItemType,
					Description: "Add an item to a grocery trip",
					Args: graphql.FieldConfigArgument{
						"tripId": &graphql.ArgumentConfig{
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
				"updateTrip": &graphql.Field{
					Type:        gql.GroceryTripType,
					Description: "Update the details about a grocery trip",
					Args: graphql.FieldConfigArgument{
						"tripId": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
						"name": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"completed": &graphql.ArgumentConfig{
							Type: graphql.Boolean,
						},
						"copyRemainingItems": &graphql.ArgumentConfig{
							Type: graphql.Boolean,
						},
					},
					Resolve: resolvers.UpdateTripResolver,
				},
			},
		},
	)

	// Define the root subscription type
	subscriptionType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Subscription",
		Fields: graphql.Fields{
			"newItem": &graphql.Field{
				Type:        gql.ItemType,
				Description: "Retrieve a new item in a trip",
				Args: graphql.FieldConfigArgument{
					"tripId": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
				},
				Resolve: subscriptions.NewItem,
			},
			// "updatedItem": &graphql.Field{
			// 	Type:        gql.ItemType,
			// 	Description: "Retrieve an update to an item",
			// 	Args: graphql.FieldConfigArgument{
			// 		"tripId": &graphql.ArgumentConfig{
			// 			Type: graphql.NewNonNull(graphql.ID),
			// 		},
			// 	},
			// 	Resolve: subscriptions.UpdatedItem,
			// },
			// "deletedItem": &graphql.Field{
			// 	Type:        gql.ItemType,
			// 	Description: "Retrieve the deletion of an item",
			// 	Args: graphql.FieldConfigArgument{
			// 		"tripId": &graphql.ArgumentConfig{
			// 			Type: graphql.NewNonNull(graphql.ID),
			// 		},
			// 	},
			// 	Resolve: subscriptions.DeletedItem,
			// },
		},
	})

	var err error
	Schema, err = graphql.NewSchema(
		graphql.SchemaConfig{
			Query:        queryType,
			Mutation:     mutationType,
			Subscription: subscriptionType,
		},
	)
	if err != nil {
		panic(err)
	}
}
