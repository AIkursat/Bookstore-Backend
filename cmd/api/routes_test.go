package main

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
)

func Test_Routes_Exist(t *testing.T) {
	TestRoutes := testApp.routes()
	chiRoutes := TestRoutes.(chi.Router)

	// these routes must exist
	routeExist(t, chiRoutes, "/users/login")
	routeExist(t, chiRoutes, "/users/logout")
	routeExist(t, chiRoutes, "/admin/users/get/{id}")
	routeExist(t, chiRoutes, "/admin/users/save")
	routeExist(t, chiRoutes, "/admin/users")
	routeExist(t, chiRoutes, "/admin/users/delete")

}

func routeExist(t *testing.T, routes chi.Router, route string) {
   // assume that the route doesn't exist
   found := false

   // scan all the registered routes
   _ = chi.Walk(routes, func(method string, foundRoute string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
	// if we find the route we're looking for, set found to true
	if route == foundRoute {
		found = true
	}
	return nil
   })

   // if we didn't find the route, make it error
   if !found{
	t.Errorf("didn't find %s in registered routes", route)
   }

}