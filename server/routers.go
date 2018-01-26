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

  /*

  HERE WAS VOVA'S ROUTES

  Route {
    "NewProfileByAdmin",
    "POST",
    "/admin/profile",
    createProfile,
  },

  Route {
    "UpdateProfileByAdmin",
    "PUT",
    "/admin/profile/{id:[0-9]+}",
    updateProfile,
  },

  Route {
    "DeleteProfileByAdmin",
    "DELETE",
    "/admin/profile/{id:[0-9]+}",
    deleteProfile,
  },

  // and so on, just add new Route structs to this array*/
  Route {
    "CreateNewPost",
    "POST",
    "/post",
    CreateNewPost,
  },

  Route {
    "DeletePost",
    "DELETE",
    "/post/{post_id}",
    DeletePost,
  },

  Route {
    "GetPostByPostID",
    "GET",
    "/post/{post_id}",
    GetPost,
  },

  Route {
    "GetPostByUserID",
    "GET",
    "/post/user/{user_id}",
    GetPostByUserID,
  },

  Route {
    "UpdatePost",
    "PUT",
    "/post/{post_id}",
    UpdatePost,
  },

  Route {
    "DescribeTablePost",
    "GET",
    "/describeTablePost",
    DescribeTablePost,
  },

}
