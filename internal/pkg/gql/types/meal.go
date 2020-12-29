package gql

import (
	"github.com/graphql-go/graphql"
)

// MealType defines the gql type for Meal
var MealType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Meal",
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
			"mealType": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"notes": &graphql.Field{
				Type: graphql.String,
			},
			"date": &graphql.Field{
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
