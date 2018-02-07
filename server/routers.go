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

func handl(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("client/index.html")
	if err != nil {
		fmt.Fprintf(w, "%s", "Error")
	} else {
		tmpl.Execute(w, r)
	}
	//fmt.Fprintf(w, "Hello Home! %s", r.URL.Path[1:]) )

}

var routes = Routes{

	Route{
		"Index",
		"GET",
		"/",
		handl,
	},

	Route{
		"NewProfileByAdmin",
		"POST",
		"/admin/profile",
		createProfile(userService),
	},

	Route{
		"UpdateProfileByAdmin",
		"PUT",
		"/admin/profile/{user_id:[0-9]+}",
		updateProfile(userService),
	},

	Route{
		"DeleteProfileByAdmin",
		"DELETE",
		"/admin/profile/{user_id:[0-9]+}",
		deleteProfile(userService),
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

	// and so on, just add new Route structs to this array
}

//Initialise services here
var userService = &controllers.UserService{mysql.DB}
