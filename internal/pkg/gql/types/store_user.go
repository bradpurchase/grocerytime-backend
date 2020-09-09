package gql

import (
	"errors"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
)

// StoreUserType defines a graphql type for StoreUser
var StoreUserType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "StoreUser",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"storeId": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"userId": &graphql.Field{
				Type: graphql.ID,
			},
			"email": &graphql.Field{
				Type: graphql.String,
			},
			"creator": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Boolean),
			},
			"active": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Boolean),
			},
			"createdAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"updatedAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"user": &graphql.Field{
				Type: UserType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					userID := p.Source.(models.StoreUser).UserID
					user := &models.User{}
					if err := db.Manager.Where("id = ?", userID).First(&user).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
						return nil, err
					}
					return user, nil
				},
			},
		},
	},
)
