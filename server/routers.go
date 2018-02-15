package server

import (
	"fmt"
	"go-team-room/controllers"
	"go-team-room/models/dao/mysql"
	"html/template"
	"net/http"

	"go-team-room/controllers/messages"

	"github.com/gorilla/mux"
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
    handler = Authorize(handler)
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

func Handl(w http.ResponseWriter, r *http.Request) {
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
    "Login",
    "POST",
    "/login",
    loginhandler,
  },

  Route{
    "Logout",
    "GET",
    "/logout",
    logout,
  },

  Route{
    "GetMessage",
    "GET",
    "/messages",
    //Test this rout by next string
    //curl -X GET "http://localhost:8080/messages?id=33&numberOfMessages=1" -H  "accept: application/json"
    messages.HandlerOfGetMessages,
  },
  Route{
    "PutMessage",
    "POST",
    "/messages",
    //Test this rout by next string
    //curl -X POST --header 'Content-Type: application/json' --header 'Accept: application/json' -d '{"message_chat_room_id": "997","message_data": {"binary_parts": [{"bin_data": null,"bin_name": null }],"text": "0 A lot of text and stupid smiles :)))))","type": "TypeOfHumMessage-UNDEFINED FOR NOW"},"message_id": "20180110155343150","message_parent_id": "","message_social_status": {"Dislike": 10,"Like": 222,"Views": 303 },"message_timestamp": "20180110155533111","message_user": {"id_sql": 13,"name_sql": "Vasya" }}' 'http://localhost:8080/messages'
    messages.HandlerOfPOSTMessages,
  },

  Route {
    "RegisterUser",
    "POST",
    "/registration",
    registerUser(userService),
  },

  Route {
    "RecoveryPass",
    "GET",
    "/recoveryPass",
    recoveryPass(userService),
  },


  // and so on, just add new Route structs to this array
}

//Initialise services here
var userService = &controllers.UserService{mysql.DB}
