package api

import (
	"log"
	"os"

	"github.com/derekdowling/jsh-api"
	"github.com/rs/cors"
	"github.com/zenazn/goji/web"
)

// BuildServer creates a new HTTP Server
func BuildServer() *web.Mux {

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{"POST", "PATCH", "GET"},
		Debug:          true,
	})

	api := jshapi.New("")
	jshapi.Logger = log.New(os.Stdout, "api", log.Ldate|log.Ltime)
	api.Use(c.Handler)

	userStorage := &UserStorage{}
	users := jshapi.NewCRUDResource("users", userStorage)

	api.Add(users)

	return api
}
