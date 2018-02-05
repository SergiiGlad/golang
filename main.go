package main

import (
  "net/http"
  "go-team-room/conf"
  "go-team-room/server"
  "fmt"
)

func main() {

  r := server.NewRouter()
  r.PathPrefix("/api-docs/").Handler(http.StripPrefix("/api-docs/", http.FileServer(http.Dir("swagger"))))
  r.PathPrefix("/dist/").Handler(http.StripPrefix("/dist/", http.FileServer(http.Dir("client/dist"))))
  http.Handle("/", r)
  http.HandleFunc("/login", server.Loginhandler)

  //r.Use(server.Authorize)

  var err = http.ListenAndServe(conf.Ip + ":" + conf.Port, nil)
  if err != nil {
    fmt.Println("Error")
  }
}
