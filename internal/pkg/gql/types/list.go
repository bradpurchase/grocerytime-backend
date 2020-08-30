package gql

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/gql/resolvers"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/grocerylist"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/trips"
	"github.com/graphql-go/graphql"
)

// ListType defines a graphql type for List
var ListType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "List",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"userId": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"creator": &graphql.Field{
				Type:    BasicUserType,
				Resolve: resolvers.BasicUserResolver,
			},
			"store": &graphql.Field{
				Type: StoreType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					db := db.FetchConnection()
					defer db.Close()

					listID := p.Source.(models.List).ID
					store := &models.Store{}
					if err := db.Where("list_id = ?", listID).First(&store).Error; err != nil {
						return nil, err
					}

					return store, nil
				},
			},
			"trip": &graphql.Field{
				Type: GroceryTripType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					db := db.FetchConnection()
					defer db.Close()

					header := p.Info.RootValue.(map[string]interface{})["Authorization"]
					user, err := auth.FetchAuthenticatedUser(db, header.(string))
					if err != nil {
						return nil, err
					}

					listID := p.Source.(models.List).ID
					trip, err := trips.RetrieveCurrentTripInList(db, listID, user.(models.User))
					if err != nil {
						return nil, err
					}
					return trip, nil
				},
			},
			"listUsers": &graphql.Field{
				Type: graphql.NewList(ListUserType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					db := db.FetchConnection()
					defer db.Close()

					listID := p.Source.(models.List).ID
					listUsers, err := grocerylist.RetrieveListUsers(db, listID)
					if err != nil {
						return nil, err
					}
					return listUsers, nil
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
