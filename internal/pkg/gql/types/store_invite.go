package gql

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/stores"
	"github.com/graphql-go/graphql"
)

// StoreInviteType defines a graphql type for Store
var StoreInviteType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "StoreInvite",
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

					// Return the store user who created this store, as this will always
					// be the user who sent the invite
					storeID := p.Source.(models.Store).ID
					user, err := stores.RetrieveStoreCreator(db, storeID)
					if err != nil {
						return nil, err
					}
					return user, nil
				},
			},
		},
	},
)
