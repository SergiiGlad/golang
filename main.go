package main

import (
	"net/http"
  "go-team-room/conf"
  "go-team-room/server"
)




func main() {
  r := server.NewRouter()

  http.Handle("/", r)

  http.ListenAndServe(conf.Ip + ":" + conf.Port, nil)
}








