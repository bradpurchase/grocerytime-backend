package gql

import (
	"errors"

	"gorm.io/gorm"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/graphql-go/graphql"
)

// ItemType defines a graphql type for Item
var ItemType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Item",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"groceryTripId": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"categoryId": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"categoryName": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					groceryTripCategoryID := p.Source.(*models.Item).CategoryID
					groceryTripCategory := &models.GroceryTripCategory{}
					if err := db.Manager.Select("store_category_id").Where("id = ?", groceryTripCategoryID).First(&groceryTripCategory).Error; err != nil {
						return nil, err
					}
					storeCategory := &models.StoreCategory{}
					if err := db.Manager.Select("name").Where("id = ?", groceryTripCategory.StoreCategoryID).First(&storeCategory).Error; err != nil {
						return nil, err
					}

					return storeCategory.Name, nil
				},
			},
			"userId": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"quantity": &graphql.Field{
				Type: graphql.Int,
			},
			"notes": &graphql.Field{
				Type: graphql.String,
			},
			"position": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"completed": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Boolean),
			},
			"user": &graphql.Field{
				Type: UserType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					userID := p.Source.(models.Item).UserID
					user := &models.User{}
					if err := db.Manager.Where("id = ?", userID).First(&user).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
						return nil, err
					}
					return user, nil
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
