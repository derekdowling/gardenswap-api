package api

import (
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/meatballhat/negroni-logrus"
	"github.com/rs/cors"
)

// API holds state for the server
type API struct {
	Logger *logrus.Logger
}

// BuildServer creates a new HTTP Server
func BuildServer() *negroni.Negroni {

	api := &API{}
	configureLogging(api)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{"POST", "PATCH", "GET"},
		Debug:          true,
	})

	server := negroni.New()
	server.Use(negronilogrus.NewMiddlewareFromLogger(api.Logger, "router"))
	server.Use(corsHandler)

	router := buildRouter(api)
	server.UseHandler(router)

	return server
}

func configureLogging(api *API) {
	api.Logger = logrus.New()
	api.Logger.Formatter = &logrus.JSONFormatter{}
}

func buildRouter(api *API) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/users", api.Register).Methods("POST").Name("register")
	router.HandleFunc("/users", api.ListUsers).Methods("GET")

	return router
}
