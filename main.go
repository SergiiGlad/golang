package main

import (
  "go-team-room/conf"
  "go-team-room/server"
  "net/http"
  "fmt"
)

func main() {
  r := server.NewRouter()
  r.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", http.FileServer(http.Dir("./swagger"))))
  r.PathPrefix("/dist/").Handler(http.StripPrefix("/dist/", http.FileServer(http.Dir("./client/dist"))))
  http.Handle("/", r)
  var err = http.ListenAndServe(conf.Ip+":"+conf.Port, nil)
  if err != nil {
    fmt.Println("Error")
  }
}
