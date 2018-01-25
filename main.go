package main

import (
	"fmt"
	"go-team-room/conf"
	"html/template"
	"net/http"

	"go-team-room/controllers/messages"
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
	http.HandleFunc("/messages", messages.HandlerOfMessages)
	http.Handle("/dist/", http.StripPrefix("/dist/", http.FileServer(http.Dir("client/dist"))))
	http.Handle("/api-docs/", http.StripPrefix("/api-docs/", http.FileServer(http.Dir("swagger"))))
	http.ListenAndServe(conf.Ip+":"+conf.Port, nil)

}
