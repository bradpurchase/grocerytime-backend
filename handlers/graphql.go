package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/gql"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

// GraphQLHandler manages all GraphQL requests
func GraphQLHandler() http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		// parse http.Request into handler.RequestOptions
		opts := handler.NewRequestOptions(request)

		rootValue := map[string]interface{}{
			"Authorization": request.Header.Get("Authorization"),
		}
		params := graphql.Params{
			Schema:         gql.Schema,
			RequestString:  opts.Query,
			VariableValues: opts.Variables,
			OperationName:  opts.OperationName,
			RootObject:     rootValue,
		}

		result := graphql.Do(params)
		payload, err := json.Marshal(result)
		if err != nil {
			log.Println("[CommandHandler] Unable to marshal JSON for publishing: ", err)
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}

		// Publish to graphql pub/sub for subscriptions here if the dataset is not nil
		//TODO clean this up
		operationName := opts.OperationName
		if operationName == "AddItemToTrip" && result.Data.(map[string]interface{})["addItemToTrip"] != nil {
			fmt.Printf("publishing message %v\n", result.Data)
			gqlPubSub.Publish("newItem", result.Data)
		}
		// if operationName == "UpdateItem" && result.Data.(map[string]interface{})["updateItem"] != nil {
		// 	fmt.Printf("publishing message %v\n", result.Data)
		// 	gqlPubSub.Publish("updatedItem", result.Data)
		// }
		// if operationName == "DeleteItem" && result.Data.(map[string]interface{})["deletedItem"] != nil {
		// 	fmt.Printf("publishing message %v\n", result.Data)
		// 	gqlPubSub.Publish("deletedItem", result.Data)
		// }

		response.WriteHeader(http.StatusOK)
		response.Header().Set("Access-Control-Allow-Origin", "*")
		response.Write(payload)
	}
}
