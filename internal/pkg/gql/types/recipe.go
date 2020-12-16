package gql

import "github.com/graphql-go/graphql"

// RecipeType defines the gql type for Recipe
var RecipeType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Recipe",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"mealType": &graphql.Field{
				Type: graphql.String,
			},
			"url": &graphql.Field{
				Type: graphql.String,
			},
			"ingredients": &graphql.Field{
				Type: RecipeIngredientType,
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
