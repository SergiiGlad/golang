package main

import (
	"fmt"
	"go-team-room/conf"
	"go-team-room/server"
	"net/http"

	"github.com/rs/cors"

	_ "go-team-room/humstat"
)


func main() {
	r := server.NewRouter()
	r.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", http.FileServer(http.Dir("./swagger"))))
	r.PathPrefix("/dist/").Handler(http.StripPrefix("/dist/", http.FileServer(http.Dir("./client/dist"))))
	r.PathPrefix("/logs/").Handler(http.StripPrefix("/logs/", http.FileServer(http.Dir("./logs"))))
	handler := cors.Default().Handler(r)
	http.Handle("/", handler)
	var err = http.ListenAndServe(conf.Ip+":"+conf.Port, nil)
	if err != nil {
		fmt.Println("Error")
	}
}
