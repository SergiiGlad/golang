package server

import (
  "net/http"
  "io/ioutil"
  "log"
  "fmt"
  "encoding/json"
  "go-team-room/models/dto"
  //"go-team-room/controllers"
  //"go-team-room/models"
)

func createUserByAdmin(w http.ResponseWriter, r *http.Request) {

  body, err := ioutil.ReadAll(r.Body)

  if err != nil {
    log.Fatal(err)
    fmt.Fprint(w, err)
  }

  user := dto.RequestUserDto{}
  err = json.Unmarshal(body, &user)

  fmt.Printf("%s", user)

  if err != nil {
    log.Fatal(err)
    fmt.Fprint(w, err)
  }

  //userDao := models.RequestUserDtoToDao(user)
  //
  //err = controllers.CreateUser(&userDao)

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
