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
				"stores": &graphql.Field{
					Type:        graphql.NewList(gql.StoreType),
					Description: "Retrieve stores for the current user",
					Resolve:     resolvers.StoresResolver,
				},
				"invitedStores": &graphql.Field{
					Type:        graphql.NewList(gql.StoreInviteType),
					Description: "Retrieve stores the current user has been invited to",
					Resolve:     resolvers.InvitedStoresResolver,
				},
				"store": &graphql.Field{
					Type:        gql.StoreType,
					Description: "Retrieve a specific store",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
					},
					Resolve: resolvers.StoreResolver,
				},
				"trips": &graphql.Field{
					Type:        graphql.NewList(gql.GroceryTripType),
					Description: "Retrieve trip history for a store",
					Args: graphql.FieldConfigArgument{
						"storeId": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
						"completed": &graphql.ArgumentConfig{
							Type: graphql.Boolean,
						},
					},
					Resolve: resolvers.GroceryTripsResolver,
				},
				"trip": &graphql.Field{
					Type:        gql.GroceryTripType,
					Description: "Retrieve a specific grocery trip within a store",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
					},
					Resolve: resolvers.GroceryTripResolver,
				},
				"itemSearch": &graphql.Field{
					Type:        gql.ItemType,
					Description: "Search for an item in the user's stores by name",
					Args: graphql.FieldConfigArgument{
						"name": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
					},
					Resolve: resolvers.ItemSearchResolver,
				},
			},
		},
	)

	// Define the root mutation type
	mutationType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Mutation",
			Fields: graphql.Fields{
				"login": &graphql.Field{
					Type:        gql.UserType,
					Description: "Retrieve an access token",
					Args: graphql.FieldConfigArgument{
						"email": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"password": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: resolvers.LoginResolver,
				},
				"token": &graphql.Field{
					Type:        gql.AuthTokenType,
					Description: "Retrieve an access token (DEPRECATED)",
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
					},
					Resolve: resolvers.TokenResolver,
				},
				"signup": &graphql.Field{
					Type:        gql.UserType,
					Description: "Create a new user account",
					Args: graphql.FieldConfigArgument{
						"email": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"password": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"name": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
					},
					Resolve: resolvers.SignupResolver,
				},
				"createStore": &graphql.Field{
					Type:        gql.StoreType,
					Description: "Create a store",
					Args: graphql.FieldConfigArgument{
						"name": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
					},
					Resolve: resolvers.CreateStoreResolver,
				},
				"deleteStore": &graphql.Field{
					Type:        gql.StoreType,
					Description: "Delete a store",
					Args: graphql.FieldConfigArgument{
						"storeId": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
					},
					Resolve: resolvers.DeleteStoreResolver,
				},
				"updateStore": &graphql.Field{
					Type:        gql.StoreType,
					Description: "Update a store",
					Args: graphql.FieldConfigArgument{
						"storeId": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
						"name": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: resolvers.UpdateStoreResolver,
				},
				"inviteToStore": &graphql.Field{
					Type:        gql.StoreUserType,
					Description: "Invite to a store via email",
					Args: graphql.FieldConfigArgument{
						"storeId": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
						"email": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
					},
					Resolve: resolvers.InviteToStoreResolver,
				},
				"joinStore": &graphql.Field{
					Type:        gql.StoreUserType,
					Description: "Removes the pending state from a pending store user",
					Args: graphql.FieldConfigArgument{
						"storeId": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
					},
					Resolve: resolvers.JoinStoreResolver,
				},
				"declineStoreInvite": &graphql.Field{
					Type:        gql.StoreUserType,
					Description: "Declines a store invitation for a user",
					Args: graphql.FieldConfigArgument{
						"storeId": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
					},
					Resolve: resolvers.DeclineStoreInviteResolver,
				},
				"leaveStore": &graphql.Field{
					Type:        gql.StoreUserType,
					Description: "Deletes the current user's store user record for the given storeID",
					Args: graphql.FieldConfigArgument{
						"storeId": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
					},
					Resolve: resolvers.LeaveStoreResolver,
				},
				"deleteItem": &graphql.Field{
					Type:        gql.ItemType,
					Description: "Remove an item from a trip",
					Args: graphql.FieldConfigArgument{
						"itemId": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
					},
					Resolve: resolvers.DeleteItemResolver,
				},
				"updateItem": &graphql.Field{
					Type:        gql.ItemType,
					Description: "Updates the properties of an item",
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
				"markItemAsCompleted": &graphql.Field{
					Type:        gql.ItemType,
					Description: "Marks item as completed by name for user (in any store)",
					Args: graphql.FieldConfigArgument{
						"name": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
					},
					Resolve: resolvers.MarkItemAsCompletedResolver,
				},
				"reorderItem": &graphql.Field{
					Type:        gql.GroceryTripType,
					Description: "Updates the order of an item in a trip",
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
						"categoryName": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: resolvers.AddItemToTrip,
				},
				"addItemsToStore": &graphql.Field{
					Type:        graphql.NewList(gql.ItemType),
					Description: "Add an array of items to a store by name. (Creates store if needed)",
					Args: graphql.FieldConfigArgument{
						"items": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(
								graphql.NewList(graphql.String),
							),
						},
						"storeName": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
					},
					Resolve: resolvers.AddItemsToStore,
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
