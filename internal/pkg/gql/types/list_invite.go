package gql

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/grocerylist"
	"github.com/graphql-go/graphql"
)

// ListInviteType defines a graphql type for List
var ListInviteType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "ListInvite",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"createdAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"updatedAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"invitingUser": &graphql.Field{
				Type:        UserType,
				Description: "Returns the user who sent the invite",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					db := db.FetchConnection()
					defer db.Close()

					// Return the list user who created this list, as this will always
					// be the user who sent the invite
					listID := p.Source.(models.List).ID
					user, err := grocerylist.RetrieveListCreator(db, listID)
					if err != nil {
						return nil, err
					}
					return user, nil
				},
			},
		},
	},
)
