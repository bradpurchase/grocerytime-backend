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
			"storeName": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					storeID := p.Source.(models.GroceryTrip).StoreID
					store := models.Store{}
					if err := db.Manager.Select("name").Where("id = ?", storeID).First(&store).Error; err != nil {
						return nil, err
					}
					return store.Name, nil
				},
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
					tripID := p.Source.(models.GroceryTrip).ID
					categories := []models.GroceryTripCategory{}
					query := db.Manager.
						Joins("INNER JOIN store_categories ON store_categories.id = grocery_trip_categories.store_category_id").
						Where("grocery_trip_categories.grocery_trip_id = ?", tripID).
						Order("grocery_trip_categories.created_at DESC").
						Find(&categories).
						Error
					if err := query; err != nil {
						return nil, err
					}
					return categories, nil
				},
			},
		},
	},
)
