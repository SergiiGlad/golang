package server

import (
  "net/http"
  "go-team-room/controllers"
  "io/ioutil"
  "go-team-room/models/dto"
  "encoding/json"
)

func Loginhandler(w http.ResponseWriter, r *http.Request) {
  body, err := ioutil.ReadAll(r.Body)

  if err != nil {
    responseError(w, err)
    return
  }

  var login dto.LoginDto

  err = json.Unmarshal(body, &login)

  if err != nil {
    responseError(w, err)
    return
  }

  user, err :=  controllers.Login(login.PhoneOrEmail, login.Password)

  if err != nil {
    responseError(w, err)
    return
  }

  var userRes dto.ResponseUserDto

  userRes = dto.UserEntityToResponseDto(user)
  userResMarsh, err := json.Marshal(userRes)

  if err != nil {
    responseError(w, err)
    return
  }

  session, err := store.Get(r, "name")

  if err != nil {
    responseError(w, err)
    return
  }

  session.Values["auth"] = "loginned"
  session.Save(r, w)
  w.Write(userResMarsh)
}
