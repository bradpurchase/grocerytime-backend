package gql

import "github.com/graphql-go/graphql"

// RecipeIngredientType defines the gql type for RecipeIngredient
var RecipeIngredientType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "RecipeIngredient",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"recipeId": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"amount": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"units": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
