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

					userID := p.Source.(models.Item).UserID
					user := &models.User{}
					if err := db.Where("id = ?", userID).First(&user).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
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
