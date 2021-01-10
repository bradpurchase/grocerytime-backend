package handlers

import (
	"encoding/json"
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
			"App-Scheme":    request.Header.Get("App-Scheme"),
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

		response.WriteHeader(http.StatusOK)
		response.Header().Set("Access-Control-Allow-Origin", "*")
		response.Write(payload)
	}
}
