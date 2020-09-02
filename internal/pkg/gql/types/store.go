package gql

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/gql/resolvers"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/stores"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/trips"
	"github.com/graphql-go/graphql"
)

// StoreType defines a graphql type for Store
var StoreType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Store",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"userId": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"creator": &graphql.Field{
				Type:    BasicUserType,
				Resolve: resolvers.BasicUserResolver,
			},
			"trip": &graphql.Field{
				Type: GroceryTripType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					db := db.FetchConnection()
					defer db.Close()

					header := p.Info.RootValue.(map[string]interface{})["Authorization"]
					user, err := auth.FetchAuthenticatedUser(db, header.(string))
					if err != nil {
						return nil, err
					}

					storeID := p.Source.(models.Store).ID
					trip, err := trips.RetrieveCurrentStoreTrip(db, storeID, user.(models.User))
					if err != nil {
						return nil, err
					}
					return trip, nil
				},
			},
			"categories": &graphql.Field{
				Type: graphql.NewList(StoreCategoryType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					db := db.FetchConnection()
					defer db.Close()

					storeID := p.Source.(models.Store).ID
					categories := []models.StoreCategory{}
					if err := db.Where("store_id = ?", storeID).Order("created_at DESC").Find(&categories).Error; err != nil {
						return nil, err
					}
					return categories, nil
				},
			},
			"users": &graphql.Field{
				Type: graphql.NewList(StoreUserType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					db := db.FetchConnection()
					defer db.Close()

					storeID := p.Source.(models.Store).ID
					storeUsers, err := stores.RetrieveStoreUsers(db, storeID)
					if err != nil {
						return nil, err
					}
					return storeUsers, nil
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
