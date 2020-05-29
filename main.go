package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/bradpurchase/grocerytime-backend/handlers"

	// Autoload env variables from .env

	_ "github.com/joho/godotenv/autoload"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"

	"github.com/gorilla/mux"
)

func main() {
	// Establish a DB factory so we can automatically run migrations etc on load
	orm := db.Factory()
	defer orm.DB.Close()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", heartbeat)
	router.Handle("/graphql", corsHandler(handlers.GraphQLHandler()))
	router.Handle("/subscriptions", handlers.WebsocketHandler())

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

func corsHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, Accept-Encoding")

		h.ServeHTTP(w, r)
	})
}
