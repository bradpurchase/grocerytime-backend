package gql

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/graphql-go/graphql"
)

var GroceryTripType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "GroceryTrip",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"storeID": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"completed": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Boolean),
			},
			"createdAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"updatedAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"categories": &graphql.Field{
				Type: graphql.NewList(GroceryTripCategoryType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					db := db.FetchConnection()
					defer db.Close()

					tripID := p.Source.(models.GroceryTrip).ID
					categories := []models.GroceryTripCategory{}
					if err := db.Where("grocery_trip_id = ?", tripID).Order("created_at DESC").Find(&categories).Error; err != nil {
						return nil, err
					}
					return categories, nil
				},
			},
		},
	},
)
