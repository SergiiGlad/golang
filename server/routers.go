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
		"/messages/",
		//HandlerOfMessages,
		messages.HandlerOfGetMessages,
	},
	Route{
		"PutMessage",
		"POST",
		"/messages/",
		//HandlerOfMessages,
		messages.HandlerOfMessages,
	},

	// and so on, just add new Route structs to this array
}

//Initialise services here
var userService = &controllers.UserService{mysql.DB}
