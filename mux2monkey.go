package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"fmt"

	"time"

	"bytes"
	"encoding/json"
)

//for routing
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
	Route{
		"Start",
		"GET",
		"/start",

		handlerStart,
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
	routePattern   string        `json:"route_pattern"`
	RemoteAddr     string        `json:"remote_addr"`
	RequestTime    time.Time     `json:"request_time"`
	RequestMethod  string        `json:"request_method"`
	Request        string        `json:"request"`
	ServerProtocol string        `json:"server_protocol"`
	Host           string        `json:"host"`
	StatusCode     int           `json:"status_code"`
	StatusText     string        `json:"status_text"`
	BodyBytesSent  int64         `json:"body_bytes_sent"`
	ElapsedTime    time.Duration `json:"elapsed_time"`
	HTTPReferrer   string        `json:"http_referrer"`
	HTTPUserAgent  string        `json:"http_user_agent"`
	RemoteUser     string        `json:"remote_user"`
	nextHandler    http.Handler
}

func (r *recHttpJson) json() ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)

	err := encoder.Encode(r)
	return buffer.Bytes(), err
}

var (
	recLogger *recHttpJson
)

func logRecord() {

	var msg string

	if jsonData, err := recLogger.json(); err != nil {
		msg = fmt.Sprintf(`{"Error": "%s"}`, err)
	} else {

		msg = string(jsonData)

	}

	fmt.Printf("Logging : %v \n", msg)
}

func recRequest(r *http.Request) {
	// do this w , r

	println("recRequest")

	recLogger.RemoteAddr = r.RemoteAddr //context.Input.IP()
	recLogger.RequestMethod = r.Method
	recLogger.ServerProtocol = r.Proto
	recLogger.Host = r.Host
	recLogger.HTTPReferrer = r.Header.Get("Referer")
	recLogger.HTTPUserAgent = r.Header.Get("User-Agent")
	recLogger.RemoteUser = r.Header.Get("Remote-User")

}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	// WriteHeader(int) is not called if our response implicitly returns 200 OK, so
	// we default to that status code.
	return &loggingResponseWriter{w, http.StatusOK}
}

func logWritter(next http.Handler, name string) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		println("logWriter")

		lrw := NewLoggingResponseWriter(w)
		next.ServeHTTP(lrw, req)

		println("come back logWriter")

		recLogger.StatusCode = lrw.statusCode
		recLogger.StatusText = http.StatusText(lrw.statusCode)
		fmt.Fprintf(w, "<-- %s %d %s\n", name, lrw.statusCode, http.StatusText(lrw.statusCode))
	})

}

func logRequest(next http.Handler, str string) http.Handler {

	//here will be done once for each route then started router

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		recLogger = new(recHttpJson)

		if recLogger == nil {
			panic("Logger record  doesn't created")
		} else {
			fmt.Println("recLogger init")
		}

		recLogger.RequestTime = time.Now()

		recRequest(req)

		fmt.Fprintf(w, "http -->  %v Route pattern :  %v %v %v %v\n", str,
			recLogger.RemoteAddr, //context.Input.IP()
			recLogger.RequestMethod,
			recLogger.ServerProtocol,
			recLogger.Host)

		next.ServeHTTP(w, req)

		println("come back log record")

		recLogger.ElapsedTime = time.Since(recLogger.RequestTime)
		logRecord()

	})

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
	w.Write([]byte("Hello from Admin \\n"))
}

func handlerStart(w http.ResponseWriter, r *http.Request) {
	println("loggin start handler")
	w.Write([]byte("Hello from start \\n"))
	http.Redirect(w, r, "/", 301)
}

func main() {

	r := NewRouter()

	panic(http.ListenAndServe("localhost:8080", r))

}
