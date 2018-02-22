package server

import (
  "fmt"
  "go-team-room/controllers"
  "go-team-room/models/dao/mysql"
  "html/template"
  "net/http"

  "github.com/gorilla/mux"
  "go-team-room/controllers/messages"
  "go-team-room/models/Amazon"
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
    //handler = Authorize(handler)
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
    "/api/admin/profile",
    createProfileByAdmin(userService),
  },

  Route{
  "UpdateProfileByAdmin",
  "PUT",
  "/api/admin/profile/{user_id:[0-9]+}",
  updateProfileByAdmin(userService),
},

  Route{
  "DeleteProfileByAdmin",
  "DELETE",
  "/api/admin/profile/{user_id:[0-9]+}",
  deleteProfileByAdmin(userService),
},

  Route {
    "CreateNewPost",
    "POST",
    "/api/post",
    CreateNewPost(Amazon.Dynamo.Db, Amazon.S3.S3API),
  },

  Route {
    "DeletePost",
    "DELETE",
    "/api/post/{post_id}",
    DeletePost(Amazon.Dynamo.Db, Amazon.S3.S3API),
  },

  Route {
    "GetPostByPostID",
    "GET",
    "/api/post/{post_id}",
    GetPost(Amazon.Dynamo.Db),
  },

  Route {
    "GetPostByUserID",
    "GET",
    "/api/post/user/{user_id}",
    GetPostByUserID(Amazon.Dynamo.Db),
  },

  Route {
    "UpdatePost",
    "PUT",
    "/api/post/{post_id}",
    UpdatePost(Amazon.Dynamo.Db),
  },

  Route {
    "GetFile",
    "GET",
    "/api/uploads/{file_link}",
    GetFileFromS3(Amazon.S3.S3API),
  },
  Route{
    "Login",
    "POST",
    "/api/login",
    loginhandler,
  },

  Route{
    "Logout",
    "GET",
    "/api/logout",
    logout,
  },

  Route{
    "GetMessage",
    "GET",
    "/api/messages",
    //Test this rout by next string
    //curl -X GET "http://localhost:8080/messages?id=33&numberOfMessages=1" -H  "accept: application/json"
    messages.HandlerOfGetMessages,
  },
  Route{
    "PutMessage",
    "POST",
    "/api/messages",
    //Test this rout by next string
    //curl -X POST --header 'Content-Type: application/json' --header 'Accept: application/json' -d '{"message_chat_room_id": "997","message_data": {"binary_parts": [{"bin_data": null,"bin_name": null }],"text": "0 A lot of text and stupid smiles :)))))","type": "TypeOfHumMessage-UNDEFINED FOR NOW"},"message_id": "20180110155343150","message_parent_id": "","message_social_status": {"Dislike": 10,"Like": 222,"Views": 303 },"message_timestamp": "20180110155533111","message_user": {"id_sql": 13,"name_sql": "Vasya" }}' 'http://localhost:8080/messages'
    messages.HandlerOfPOSTMessages,
  },

  Route{
    "RegisterUser",
    "POST",
    "/api/registration",
    registerUser(userService, emailService),
  },

  Route{
    "RecoveryPass",
    "GET",
    "/api/recoveryPass",
    recoveryPass(userService, emailService),
  },

  Route{
    "ConfirmAccount",
    "GET",
    "/api/confirm/email/{token}",
    ConfirmAccount(tokenService),
  },

  Route{
    "GetUserFriends",
    "GET",
    "/api/profile/{user_id:[0-9]+}/friends",
    getFriends(friendService),
  },

  Route{
    "GetUsersWithRequests",
    "GET",
    "/api/profile/{user_id:[0-9]+}/friends/requests",
    getUsersWithRequests(friendService),
  },

  Route{
    "NewFriendRequest",
    "POST",
    "/api/friend",
    newFriendRequest(friendService),
  },

  Route{
    "ReplyToFriendRequest",
    "PUT",
    "/api/friend",
    replyToFriendRequest(friendService),
  },

  Route{
    "DeleteFriend",
    "DELETE",
    "/api/friend",
    deleteFriendship(friendService),
  },
  Route{
    "GetMessage",
    "GET",
    "/api/messages",
    //Test this rout by next string
    //curl -X GET "http://localhost:8080/messages?id=33&numberOfMessages=1" -H  "accept: application/json"
    messages.HandlerOfGetMessages,
  },
  Route{
    "PutMessage",
    "POST",
    "/api/messages",
    //Test this rout by next string
    //curl -X POST --header 'Content-Type: application/json' --header 'Accept: application/json' -d '{"message_chat_room_id": "997","message_data": {"binary_parts": [{"bin_data": null,"bin_name": null }],"text": "0 A lot of text and stupid smiles :)))))","type": "TypeOfHumMessage-UNDEFINED FOR NOW"},"message_id": "20180110155343150","message_parent_id": "","message_social_status": {"Dislike": 10,"Like": 222,"Views": 303 },"message_timestamp": "20180110155533111","message_user": {"id_sql": 13,"name_sql": "Vasya" }}' 'http://localhost:8080/messages'
    messages.HandlerOfPOSTMessages,
  },

  Route{
    "UploadAvatar",
    "PUT",
    "/api/profile/{user_id}/avatar",
    UploadAvatar(userService, Amazon.S3.S3API),
  },

  Route{
    "DeleteAvatar",
    "DELETE",
    "/api/profile/{user_id}/avatar",
    DeleteAvatar(userService, Amazon.S3.S3API),
  },

  Route{
    "GetProfile",
    "GET",
    "/api/profile/{user_id}",
    GetProfile(userService),
  },

  // and so on, just add new Route structs to this array
}

//Initialise services here
var emailService = &controllers.UserEmailService{
  &controllers.HermesEmailBodyGenerator{},
  &controllers.DefaultEmailSend{},
  &controllers.TokenService{
    mysql.UserDao,
    mysql.TokenDao,
  },
}
var tokenService = &controllers.TokenService{
  mysql.UserDao,
  mysql.TokenDao,
}
var friendService = &controllers.FriendService{mysql.FriendshipDao, mysql.UserDao}
var userService = &controllers.UserService{
  friendService,
  mysql.PasswordDao,
  mysql.UserDao,
}
