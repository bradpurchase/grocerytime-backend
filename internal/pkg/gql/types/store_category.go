package gql

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
)

var StoreCategoryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "StoreCategory",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"storeId": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"categoryId": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
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
			"category": &graphql.Field{
				Type: CategoryType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					db := db.FetchConnection()
					defer db.Close()

					storeCategoryID := p.Source.(models.StoreCategory).CategoryID
					category := &models.Category{}
					if err := db.Where("id = ?", storeCategoryID).First(&category).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
						return nil, err
					}
					return category, nil
				},
			},
		},
	},
)
