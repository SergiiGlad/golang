package main

import (
	"net/http"
  "go-team-room/conf"
  "github.com/gorilla/mux"
)

func main() {
  r := mux.NewRouter()
  http.ListenAndServe(conf.Port, r)
}
