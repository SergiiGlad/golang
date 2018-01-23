package server

import (
  "github.com/gorilla/mux"
  "net/http"
)

type Route struct {
  Name        string
  Method      string
  Pattern     string
  HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
  router := mux.NewRouter().StrictSlash(true)
  for _, route := range routes {
    var handler http.Handler
    handler = route.HandlerFunc
    //handler = middleware.Logger(handler, route.Name)
    //handler = middleware.Auth(handler)
    // ....
    // and so on

    router.
      Methods(route.Method).
      Path(route.Pattern).
      Name(route.Name).
      Handler(handler)
  }

  return router
}

var routes = Routes{

  Route {
    "NewProfileByAdmin",
    "POST",
    "/admin/profile",
    createUserByAdmin,
  },

  // and so on, just add new Route structs to this array
}
