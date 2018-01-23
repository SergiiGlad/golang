package main

import (
  "net/http"
  "go-team-room/conf"
  "go-team-room/server"
  "fmt"
  "html/template"
  "encoding/json"
  "go-team-room/models/dao"
)

func handl(w http.ResponseWriter, r *http.Request) {
  tmpl, err := template.ParseFiles("client/index.html")
  if err != nil {
    fmt.Fprintf(w, "%s", "Error")
  } else {
    tmpl.Execute(w, r)
  }
  //fmt.Fprintf(w, "Hello Home! %s", r.URL.Path[1:]) )

}

func main() {

  user := dao.User {
    0,
    "email@gmail.com",
    "Vova",
    "Polischuk",
    "+380509787333",
    "123123rock",
    "user",
    "active",
    "",
  }

  js, _ := json.Marshal(user)

  fmt.Printf("%s", js)

  r := server.NewRouter()
  //http.HandleFunc("/", handl)
  http.Handle("/", r)
  http.Handle("/api-docs/", http.StripPrefix("/api-docs/", http.FileServer(http.Dir("swagger"))))
  http.Handle("/dist/", http.StripPrefix("/dist/", http.FileServer(http.Dir("client/dist"))))

  http.ListenAndServe(conf.Ip+":"+conf.Port, nil)
}
