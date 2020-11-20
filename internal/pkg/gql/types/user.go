package gql

import (
	"errors"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/graphql-go/graphql"
	uuid "github.com/satori/go.uuid"
)

// UserType defines a graphql type for User
var UserType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"email": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"passwordResetToken": &graphql.Field{
				Type: graphql.String,
			},
			"passwordResetTokenExpiry": &graphql.Field{
				Type: graphql.DateTime,
			},
			"createdAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"updatedAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"defaultStoreId": &graphql.Field{
				Type: graphql.ID,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					userID := p.Source.(models.User).ID
					//TODO add index for this query?
					var storeIDs []uuid.UUID
					query := db.Manager.
						Model(&models.Store{}).
						Select("stores.id").
						Joins("INNER JOIN store_users ON store_users.store_id = stores.id").
						Joins("INNER JOIN store_user_preferences ON store_user_preferences.store_user_id = store_users.id").
						Where("store_users.user_id = ?", userID).
						Where("store_user_preferences.default_store = ?", true).
						Pluck("stores.id", &storeIDs).
						Error
					if err := query; err != nil {
						return nil, errors.New("default store id not found for user")
					}
					if len(storeIDs) > 0 {
						return storeIDs[0], nil
					}
					return nil, nil
				},
			},
			"accessToken": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					userID := p.Source.(*models.User).ID
					var authToken models.AuthToken
					query := db.Manager.
						Select("access_token").
						Where("user_id = ?", userID).
						Last(&authToken).
						Error
					if err := query; err != nil {
						return nil, errors.New("token not found for user")
					}
					return authToken.AccessToken, nil
				},
			},
			// DEPRECATED in favour of accessToken
			"token": &graphql.Field{
				Type: AuthTokenType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					userID := p.Source.(*models.User).ID
					authToken := &models.AuthToken{}
					if err := db.Manager.Where("user_id = ?", userID).Last(&authToken).Error; err != nil {
						return nil, errors.New("token not found for user")
					}
					return authToken, nil
				},
			},
		},
	},
)

// BasicUserType is similar to UserType except it does not include fields for anything
// except basic user info. Used in sharableStore query
var BasicUserType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "BasicUser",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"email": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"firstName": &graphql.Field{
				Type: graphql.String,
			},
			"lastName": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
