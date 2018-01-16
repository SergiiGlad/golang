package main

import (
	"fmt"
	"net/http"
  "html/template"
)

func handler(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("client/index.html")
    if err != nil {
      fmt.Fprintf(w, "%s", "Error")
    } else {
      tmpl.Execute(w, r)
    }
    //fmt.Fprintf(w, "Hello Home! %s", r.URL.Path[1:]) )

}

func main() {
	http.HandleFunc("/", handler)
  http.Handle("/dist/", http.StripPrefix("/dist/", http.FileServer(http.Dir("client/dist"))))
	http.ListenAndServe(":8080", nil)
}
