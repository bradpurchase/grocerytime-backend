package gql

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/trips"
	"github.com/graphql-go/graphql"
)

var GroceryTripCategoryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "GroceryTripCategory",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"storeCategory": &graphql.Field{
				Type: StoreCategoryType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					db := db.FetchConnection()
					defer db.Close()

					storeCategoryID := p.Source.(models.GroceryTripCategory).StoreCategoryID
					storeCategory := models.StoreCategory{}
					if err := db.Select("id, name").Where("id = ?", storeCategoryID).First(&storeCategory).Error; err != nil {
						return nil, err
					}
					return storeCategory, nil
				},
			},
			"items": &graphql.Field{
				Type: graphql.NewList(ItemType),
				Args: graphql.FieldConfigArgument{
					"filter": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					db := db.FetchConnection()
					defer db.Close()

					tripID := p.Source.(models.GroceryTripCategory).GroceryTripID
					categoryID := p.Source.(models.GroceryTripCategory).ID
					items, err := trips.RetrieveItemsInCategory(db, tripID, categoryID)
					if err != nil {
						return nil, err
					}
					return items, nil
				},
			},
		},
	},
)
