package server

import (
  "github.com/gorilla/mux"
  "net/http"
  "fmt"
  "html/template"
  "go-team-room/controllers"
  "go-team-room/models/dao/mysql"
)

/*
Route describes request routing (it defines what HandlerFunc
should be called for incoming request). Its fields are used to
register a new gorilla route with a matcher for HTTP methods
and the URL path.
*/
type Route struct {
  Name        string
  Method      string
  Pattern     string
  HandlerFunc http.HandlerFunc
}

type Routes []Route


//NewRouter creates new mux.Router to handle incoming requests
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

func handl(w http.ResponseWriter, r *http.Request) {
  tmpl, err := template.ParseFiles("client/index.html")
  if err != nil {
    fmt.Fprintf(w, "%s", "Error")
  } else {
    tmpl.Execute(w, r)
  }
  //fmt.Fprintf(w, "Hello Home! %s", r.URL.Path[1:]) )

	log.Info(reqtoLog(r))
}

var routes = Routes{

  Route {
    "Index",
    "GET",
    "/",
    handl,
  },

  Route {
    "NewProfileByAdmin",
    "POST",
    "/admin/profile",
    createProfileByAdmin(userService),
  },

  Route {
    "UpdateProfileByAdmin",
    "PUT",
    "/admin/profile/{user_id:[0-9]+}",
    updateProfileByAdmin(userService),
  },

  Route {
    "DeleteProfileByAdmin",
    "DELETE",
    "/admin/profile/{user_id:[0-9]+}",
    deleteProfileByAdmin(userService),
  },

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
    "GetFile",
    "GET",
    "/uploads/{file_link}",
    GetFileFromS3,
  },
  Route {
    "RegisterUser",
    "POST",
    "/registration",
    registerUser(userService),
  },

  // and so on, just add new Route structs to this array
}

//Initialise services here
var userService = &controllers.UserService{mysql.DB}
