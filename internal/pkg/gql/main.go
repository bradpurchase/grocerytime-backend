package gql

import (
	"github.com/graphql-go/graphql"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/gql/resolvers"
	gql "github.com/bradpurchase/grocerytime-backend/internal/pkg/gql/types"
)

var Schema graphql.Schema

func init() {
	// Define the root query type
	queryType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "RootQuery",
			Fields: graphql.Fields{
				"me": &graphql.Field{
					Type:        gql.UserType,
					Description: "Retrieve the current user",
					Resolve:     resolvers.AuthenticatedUserResolver,
				},
			},
		},
	)

	// Define the root mutation type
	mutationType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Mutation",
			Fields: graphql.Fields{
				"signup": &graphql.Field{
					Type:        gql.AuthTokenType,
					Description: "Create a new user account",
					Args: graphql.FieldConfigArgument{
						"email": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"password": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
					},
					Resolve: resolvers.SignupResolver,
				},
				"token": &graphql.Field{
					Type:        gql.AuthTokenType,
					Description: "Retrieve an access token",
					Args: graphql.FieldConfigArgument{
						"grantType": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"email": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"password": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"refreshToken": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: resolvers.TokenResolver,
				},
			},
		},
	)

	var err error
	Schema, err = graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    queryType,
			Mutation: mutationType,
		},
	)
	if err != nil {
		panic(err)
	}
}
