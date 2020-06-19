package gql

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
)

// ListUserType defines a graphql type for ListUser
var ListUserType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "ListUser",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"listId": &graphql.Field{
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
					db := db.FetchConnection()
					defer db.Close()

					userID := p.Source.(models.ListUser).UserID
					user := &models.User{}
					if err := db.Where("id = ?", userID).First(&user).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
						return nil, err
					}
					return user, nil
				},
			},
		},
	},
)
