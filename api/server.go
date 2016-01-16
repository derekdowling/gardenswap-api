package api

import (
	"log"
	"os"

	"goji.io/pat"

	"github.com/derekdowling/go-stdlogger"
	"github.com/derekdowling/jsh-api"
	"github.com/rs/cors"
)

// BuildServer creates a new HTTP Server
func BuildServer(debug bool) *jshapi.API {

	// setup logging
	logger := buildLogger()
	jshapi.Logger = logger

	api := jshapi.New("", debug)

	// set middleware
	api.Use(buildCORS(debug).Handler)

	// /users Routes
	userAPI := &UserAPI{Logger: logger}
	users := jshapi.NewCRUDResource("users", userAPI)
	api.Add(users)
	api.Mux.HandleFuncC(pat.Post("/register"), userAPI.Register)

	return api
}

func buildCORS(debug bool) *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{"POST", "PATCH", "GET", "DELETE"},
		Debug:          debug,
	})
}

func buildLogger() std.Logger {
	return log.New(os.Stderr, "gardenswap: ", log.LstdFlags)
}
