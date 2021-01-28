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
				"passwordReset": &graphql.Field{
					Type:        gql.UserType,
					Description: "Retrieve information about a password reset",
					Args: graphql.FieldConfigArgument{
						"token": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
					},
					Resolve: resolvers.PasswordResetResolver,
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
				"storeUserPrefs": &graphql.Field{
					Type:        gql.StoreUserPreferenceType,
					Description: "Retrieves the current user's preferences for a store",
					Args: graphql.FieldConfigArgument{
						"storeId": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
					},
					Resolve: resolvers.StoreUserPrefsResolver,
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
				"recipes": &graphql.Field{
					Type:        graphql.NewList(gql.RecipeType),
					Description: "Retrieve recipes added by the current user",
					Resolve:     resolvers.RecipesResolver,
				},
				"recipe": &graphql.Field{
					Type:        gql.RecipeType,
					Description: "Retrieve recipe",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
					},
					Resolve: resolvers.RecipeResolver,
				},
				"meals": &graphql.Field{
					Type:        graphql.NewList(gql.MealType),
					Description: "Retrieve planned meals for the current user within the provided time period",
					Args: graphql.FieldConfigArgument{
						"year": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
						"weekNumber": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
					},
					Resolve: resolvers.MealsResolver,
				},
				"meal": &graphql.Field{
					Type:        gql.MealType,
					Description: "Retrieve meal",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
					},
					Resolve: resolvers.MealResolver,
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
						"deviceName": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: resolvers.LoginResolver,
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
						"deviceName": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: resolvers.SignupResolver,
				},
				"forgotPassword": &graphql.Field{
					Type:        gql.UserType,
					Description: "Sends an email to a user to set a new password",
					Args: graphql.FieldConfigArgument{
						"email": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
					},
					Resolve: resolvers.ForgotPasswordResolver,
				},
				"resetPassword": &graphql.Field{
					Type:        gql.UserType,
					Description: "Updates a password for a user",
					Args: graphql.FieldConfigArgument{
						"password": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"token": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
					},
					Resolve: resolvers.ResetPasswordResolver,
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
				"updateStoreUserPrefs": &graphql.Field{
					Type:        gql.StoreUserPreferenceType,
					Description: "Updates the current user's preferences for a store",
					Args: graphql.FieldConfigArgument{
						"storeId": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
						"defaultStore": &graphql.ArgumentConfig{
							Type: graphql.Boolean,
						},
						"notifications": &graphql.ArgumentConfig{
							Type: graphql.Boolean,
						},
					},
					Resolve: resolvers.UpdateStoreUserPrefsResolver,
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
						"notes": &graphql.ArgumentConfig{
							Type: graphql.String,
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
							Type: graphql.String,
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
						"newTripName": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: resolvers.UpdateTripResolver,
				},
				"addDevice": &graphql.Field{
					Type:        gql.DeviceType,
					Description: "Stores a device token for the current user for push notifications",
					Args: graphql.FieldConfigArgument{
						"token": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
					},
					Resolve: resolvers.AddDeviceResolver,
				},
				"notifyTripUpdatedItemsAdded": &graphql.Field{
					Type:        graphql.Boolean,
					Description: "Notifies store users about a trip being updated after items were added by the current user",
					Args: graphql.FieldConfigArgument{
						"storeId": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
						"numItemsAdded": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.Int),
						},
					},
					Resolve: resolvers.NotifyTripUpdatedItemsAddedResolver,
				},
				"createRecipe": &graphql.Field{
					Type:        gql.RecipeType,
					Description: "Creates a recipe",
					Args: graphql.FieldConfigArgument{
						"name": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"mealType": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"url": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"ingredients": &graphql.ArgumentConfig{
							Type: graphql.NewList(graphql.NewInputObject(
								graphql.InputObjectConfig{
									Name: "RecipeIngredientInput",
									Fields: graphql.InputObjectConfigFieldMap{
										"name": &graphql.InputObjectFieldConfig{
											Type: graphql.String,
										},
										"amount": &graphql.InputObjectFieldConfig{
											Type: graphql.Float,
										},
										"unit": &graphql.InputObjectFieldConfig{
											Type: graphql.String,
										},
										"notes": &graphql.InputObjectFieldConfig{
											Type: graphql.String,
										},
									},
								},
							)),
						},
					},
					Resolve: resolvers.CreateRecipeResolver,
				},
				"planMeal": &graphql.Field{
					Type:        gql.MealType,
					Description: "Creates a planned meal of a recipe",
					Args: graphql.FieldConfigArgument{
						"recipeId": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
						"storeId": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
						"name": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"mealType": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"notes": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"servings": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
						"date": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"items": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.NewList(graphql.NewInputObject(
								graphql.InputObjectConfig{
									Name: "MealItemInput",
									Fields: graphql.InputObjectConfigFieldMap{
										"name": &graphql.InputObjectFieldConfig{
											Type: graphql.String,
										},
										"quantity": &graphql.InputObjectFieldConfig{
											Type: graphql.Int,
										},
									},
								},
							))),
						},
					},
					Resolve: resolvers.PlanMealResolver,
				},
				"deleteMeal": &graphql.Field{
					Type:        gql.MealType,
					Description: "Deletes a planned meal for the current user and anyone the meal is shared with",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
					},
					Resolve: resolvers.DeleteMealResolver,
				},
				"updateMeal": &graphql.Field{
					Type:        gql.MealType,
					Description: "Updates a planned meal for the current user and anyone the meal is shared with",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
						"name": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"mealType": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"notes": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"servings": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
						"date": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: resolvers.UpdateMealResolver,
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
