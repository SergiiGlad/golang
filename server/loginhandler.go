package server

import (
  "net/http"
  "go-team-room/controllers"
  "io/ioutil"
  "go-team-room/models/dto"
  "encoding/json"
  "github.com/gorilla/sessions"
)

func loginhandler(w http.ResponseWriter, r *http.Request) {
  body, err := ioutil.ReadAll(r.Body)

  if err != nil {
    responseError(w, err, http.StatusBadRequest)
    return
  }

  var login dto.LoginDto

  err = json.Unmarshal(body, &login)

  if err != nil {
    responseError(w, err, http.StatusBadRequest)
    return
  }

  session, err := store.Get(r, "name")

  if err != nil {
    responseError(w, err, http.StatusBadRequest)
    return
  }

  user, err := controllers.Login(login.PhoneOrEmail, login.Password)

  if err != nil {
    responseError(w, err, http.StatusForbidden)
    //session.AddFlash(errors.New("Wrong credentials"), "_errors")
    //session.Save(r, w)
    return
  }

  var userRes dto.ResponseUserDto

  userRes = dto.UserEntityToResponseDto(user)
  userResMarsh, err := json.Marshal(userRes)

  if err != nil {
    responseError(w, err, http.StatusBadRequest)
    return
  }

  session.Values["loginned"] = true
  session.Values["user_ID"] = user.ID
  session.Values["role"] = user.Role

  store.Options = &sessions.Options {
    MaxAge:   24*60*60,
    Secure:   true,
    HttpOnly: true,
      }

  session.Save(r, w)
  w.Write(userResMarsh)
}
