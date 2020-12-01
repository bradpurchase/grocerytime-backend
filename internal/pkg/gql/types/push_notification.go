package gql

import "github.com/graphql-go/graphql"

// PushNotificationType defines a graphql type for Device
var PushNotificationType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "PushNotification",
		Fields: graphql.Fields{
			"apnsId": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
	},
)
