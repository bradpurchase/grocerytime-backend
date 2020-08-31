package gql

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
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
			"userId": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"quantity": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
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
					db := db.FetchConnection()
					defer db.Close()

					userID := p.Source.(models.Item).UserID
					user := &models.User{}
					if err := db.Where("id = ?", userID).First(&user).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
						return nil, err
					}
					return user, nil
				},
			},
			"category": &graphql.Field{
				Type: StoreCategoryType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					db := db.FetchConnection()
					defer db.Close()

					categoryID := p.Source.(models.Item).CategoryID
					category := &models.StoreCategory{}
					if err := db.Where("id = ?", categoryID).First(&category).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
						return nil, err
					}
					return category, nil
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
