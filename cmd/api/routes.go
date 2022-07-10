package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler{
    mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"}, // Audience can enter with both http and https methods.
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Rest API methods will be accepted by service
		AllowedHeaders: []string{"Accept", "Auuthorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: true,
		MaxAge: 300, 
	}))


	mux.Get("/users/login", app.Login)
	
	mux.Post("/users/login", app.Login)
    
	return mux
}