package main

import (
  //"go-team-room/conf"
  //"go-team-room/server"
  //"net/http"
  //"fmt"
  "go-team-room/controllers"
  "go-team-room/models/dto"
)

func main() {
  em := dto.Email{
    To: "arsenzhd@gmail.com",
    Body: "Hello world",
    Subject: "Hello !!",
      }
 de := controllers.DefaultEmailSend{}
 e := de.SendEmail(em)
 e.Error();

	//r := server.NewRouter()
	//r.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", http.FileServer(http.Dir("./swagger"))))
	//r.PathPrefix("/dist/").Handler(http.StripPrefix("/dist/", http.FileServer(http.Dir("./client/dist"))))
	//r.PathPrefix("/logs/").Handler(http.StripPrefix("/logs/", http.FileServer(http.Dir("./logs"))))
	//http.Handle("/", r)
	//var err = http.ListenAndServe(conf.Ip+":"+conf.Port, nil)
	//if err != nil {
	//	fmt.Println("Error")
	//}
}
