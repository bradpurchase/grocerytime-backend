package gql

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
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
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"categories": &graphql.Field{
				Type: graphql.NewList(StoreCategoryType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					db := db.FetchConnection()
					defer db.Close()

					storeID := p.Source.(*models.Store).ID
					categories := []models.StoreCategory{}
					if err := db.Where("store_id = ?", storeID).Order("created_at DESC").Find(&categories).Error; err != nil {
						return nil, err
					}
					return categories, nil
				},
			},
			"createdAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"updatedAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"deletedAt": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	},
)
