package main

import (
	"log"
	"net/http"

	"github.com/derekdowling/gardenswap-api/api"
)

func main() {
	server := api.BuildServer(true)

	address := "0.0.0.0:8080"
	log.Printf("Server started, listening at %s", address)
	log.Fatal(http.ListenAndServe(address, server.Mux))
}
