package main

import (
	"Bookstore-Backend/internal/data"
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

	mux.Get("/users/all", func(w http.ResponseWriter, r *http.Request){
		var users data.User // We created the users
		all, err := users.GetAll() // from models
		if err != nil{
			app.errorLog.Println(err)
			return
		}

		app.writeJSON(w, http.StatusOK, all)

	})

	mux.Get("/users/add", func(w http.ResponseWriter, r *http.Request){
		var u = data.User{
			Email: "you@there.com",
			FirstName: "You",
			LastName: "There",
			Password: "password",
		}
		app.infoLog.Println("Adding user...")


		id, err := app.models.User.Insert(u)
		if err != nil{
			app.errorLog.Println(err)
			app.errorJSON(w, err, http.StatusForbidden)
			return
		}

		// we will get the user

		app.infoLog.Println("Got back id of", id)

	    // we will get the user from db that's why we created newUser and ignored the error	

		newUser, _ :=  app.models.User.GetOne(id)
		app.writeJSON(w, http.StatusOK, newUser)
	})
    
	return mux
}