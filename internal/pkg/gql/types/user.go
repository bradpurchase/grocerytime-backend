package gql

import (
	"errors"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/graphql-go/graphql"
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
			"token": &graphql.Field{
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
