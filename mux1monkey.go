package main

import (
	"net/http"
	"github.com/gorilla/mux"

	"fmt"


	"time"


)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{

	Route{
		"Index",
		"GET",
		"/",

		handlerMain,
	},
	Route{
		"Admin",
		"GET",
		"/admin",

		handlerAdmin,
	},
}


type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// struct for holding access log data.
type recHttpJson struct {
	RemoteAddr     string        `json:"remote_addr"`
	RequestTime    time.Time     `json:"request_time"`
	RequestMethod  string        `json:"request_method"`
	Request        string        `json:"request"`
	ServerProtocol string        `json:"server_protocol"`
	Host           string        `json:"host"`
	Status         int           `json:"status"`
	BodyBytesSent  int64         `json:"body_bytes_sent"`
	ElapsedTime    time.Duration `json:"elapsed_time"`
	HTTPReferrer   string        `json:"http_referrer"`
	HTTPUserAgent  string        `json:"http_user_agent"`
	RemoteUser     string        `json:"remote_user"`
}

type recHttpText struct {
	routePattern      string
	RemoteAddr     string
	RequestMethod  string
	RequestTime    time.Time
	Status         int
	ElapsedTime    time.Duration
	ServerProtocol string
	Host           string
	HTTPReferrer   string
	HTTPUserAgent  string
	RemoteUser     string
	nextHandler    http.Handler
}


var _ http.Handler = &recHttpText{}

func (lh *recHttpText) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// do this w , r


	lh.RequestTime = time.Now()
	lh.RemoteAddr = r.RemoteAddr //context.Input.IP()
	lh.RequestMethod = r.Method
	lh.ServerProtocol = r.Proto
	lh.Host = r.Host
	lh.HTTPReferrer = r.Header.Get("Referer")
	lh.HTTPUserAgent = r.Header.Get("User-Agent")
	lh.RemoteUser = r.Header.Get("Remote-User")


	fmt.Fprintf(w, "http -->  %v Route pattern :  %v \n", lh.routePattern, lh)


	lh.nextHandler.ServeHTTP(w, r)

	lh.ElapsedTime = time.Since( lh.RequestTime )

	//formatter := time.RFC1123
	fmt.Fprintf(w, " Duration %t ", lh.ElapsedTime)

	println("come back from lof writer")

}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	// WriteHeader(int) is not called if our response implicitly returns 200 OK, so
	// we default to that status code.
	return &loggingResponseWriter{w, http.StatusOK}
}

func logWritter(next http.Handler, name string) http.Handler {

	println("logWriter")

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		lrw := NewLoggingResponseWriter(w)
		next.ServeHTTP(lrw, req)

		statusCode := lrw.statusCode
		fmt.Fprintf(w,"<-- %s %d %s", name, statusCode, http.StatusText(statusCode))
	})

}

func logRequest(next http.Handler, str string) http.Handler {

	la := &recHttpText{
		nextHandler: next,
		routePattern:   str,
	}

	println("logHttpRequest")

	return la

}

//NewRouter creates new mux.Router to handle incoming requests
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler



		handler = route.HandlerFunc
		handler = logWritter(handler, route.Name)
		handler = logRequest(handler, route.Pattern)
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

func handlerMain(w http.ResponseWriter, r *http.Request) {

	println("loggin root handler")
	w.Write([]byte("Hello from root dir "))
}

func handlerAdmin(w http.ResponseWriter, r *http.Request) {
	println("loggin admin handler")
	w.Write([]byte("Hello from Admin "))
}

func main() {

	r := NewRouter()

	panic(http.ListenAndServe(":8080", r))

}
