package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler{
    mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)


	mux.Get("/users/login", app.Login)
	
	mux.Post("/users/login", app.Login)
    
	return mux
}