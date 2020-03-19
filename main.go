package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	// Autoload env variables from .env
	_ "github.com/joho/godotenv/autoload"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/gql"

	"github.com/gorilla/mux"
	"github.com/graphql-go/handler"
)

func main() {
	// Establish a DB factory so we can automatically run migrations etc on load
	orm := db.Factory()
	defer orm.DB.Close()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", heartbeat)
	router.Handle("/graphql", handler.New(&handler.Config{
		Schema:   &gql.Schema,
		Pretty:   true,
		GraphiQL: true,
		RootObjectFn: func(ctx context.Context, r *http.Request) map[string]interface{} {
			return map[string]interface{}{
				"Authorization": r.Header.Get("Authorization"),
			}
		},
	}))

	port := os.Getenv("PORT")
	log.Println("[main] ...Listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

type heartbeatResponse struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
}

func heartbeat(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(heartbeatResponse{Status: "OK", Code: 200})
}
