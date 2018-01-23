package server

import (
  "net/http"
  "io/ioutil"
  "log"
  "fmt"
  "encoding/json"
  "go-team-room/controllers"
  //"go-team-room/models/dao"
  "go-team-room/models/dto"
)

func createUserByAdmin(uservice controllers.UserService) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    fmt.Println("/admin/profile")

    body, err := ioutil.ReadAll(r.Body)

    if err != nil {
      log.Fatal(err)
      fmt.Fprint(w, err)
    }

    user := dto.UserDto{}
    err = json.Unmarshal(body, &user)

    fmt.Printf("%s", user)

    if err != nil {
      log.Fatal(err)
      fmt.Fprint(w, err)
    }

    //err = uservice.CreateUser(&user)

    if err != nil {
      log.Fatal(err)
      fmt.Fprint(w, err)
    }

    body, err = json.Marshal(user)
    _, err = w.Write(body)

    if err != nil {
      log.Fatal(err)
      fmt.Fprint(w, err)
    }
  }
}
