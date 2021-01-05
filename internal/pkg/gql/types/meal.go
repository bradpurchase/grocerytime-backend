package gql

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/graphql-go/graphql"
	uuid "github.com/satori/go.uuid"
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
			"servings": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"notes": &graphql.Field{
				Type: graphql.String,
			},
			"date": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"users": &graphql.Field{
				Type: graphql.NewList(UserType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					mealUsers := p.Source.(models.Meal).Users
					var userIDs []uuid.UUID
					for i := range mealUsers {
						userIDs = append(userIDs, mealUsers[i].UserID)
					}
					var users []models.User
					if err := db.Manager.Where("id IN (?)", userIDs).Find(&users).Error; err != nil {
						return nil, err
					}
					return users, nil
				},
			},
			"recipe": &graphql.Field{
				Type: RecipeType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var recipe models.Recipe
					recipeID := p.Source.(models.Meal).RecipeID
					query := db.Manager.
						Preload("Ingredients").
						Where("id = ?", recipeID).
						Last(&recipe).
						Error
					if err := query; err != nil {
						return nil, err
					}
					return recipe, nil
				},
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
