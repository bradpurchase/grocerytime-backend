package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/bradpurchase/grocerytime-backend/handlers"

	// Autoload env variables from .env

	_ "github.com/joho/godotenv/autoload"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"

	"github.com/gorilla/mux"
)

func main() {
	// Establish a DB factory so we can automatically run migrations etc on load
	db.Factory()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", heartbeat)

	// Debugging (pprof)
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profline", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	router.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
	router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	router.Handle("/debug/pprof/block", pprof.Handler("block"))

	router.Handle("/graphql", corsHandler(handlers.GraphQLHandler()))

	port := os.Getenv("PORT")
	log.Println("[main] ⚡️...Listening on port " + port)
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
