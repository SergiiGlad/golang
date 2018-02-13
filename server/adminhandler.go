package server

import (
  "net/http"
  "encoding/json"
  "go-team-room/controllers"
  "github.com/gorilla/mux"
  "strconv"
  "go-team-room/conf"
)

// Get instance of logger (Formatter, Hook，Level, Output ).
// If you want to use only your log message  It will need use own call logs example
var log = conf.GetLog()

//createProfileByAdmin is HandlerFunc wrapper. It accepts types that implement UserServiceInterface.
// This function can be called to create new user with any role (dao.Role) type. Use it only for
// admin role use cases.
func createProfileByAdmin(service controllers.UserServiceInterface) http.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request) {
    reqUserDto, err := userDtoFromReq(r)
    if err != nil {
      responseError(w, err, http.StatusForbidden)
      return
    }

    respUserDto, err := service.CreateUser(&reqUserDto)
    if err != nil {
      responseError(w, err, http.StatusForbidden)
      return
    }

    respBody, err := json.Marshal(respUserDto)
    _, err = w.Write(respBody)
    if err != nil {
      responseError(w, err, http.StatusForbidden)
      return
    }
  }
}

//updateProfileByAdmin is HandlerFunc with similar signature as createProfileByAdmin. It can be used
//to modify user. Use it only for admin role use cases.
func updateProfileByAdmin(service controllers.UserServiceInterface) http.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request) {

    userDto, err := userDtoFromReq(r)
    if err != nil {
      responseError(w, err, http.StatusForbidden)
      return
    }

    idStr := mux.Vars(r)["user_id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
      responseError(w, err, http.StatusForbidden)
      return
    }

    respUserDto, err := service.UpdateUser(int64(id), &userDto)
    if err != nil {
      responseError(w, err, http.StatusForbidden)
      return
    }

    respBody, err := json.Marshal(respUserDto)
    _, err = w.Write(respBody)
    if err != nil {
      responseError(w, err, http.StatusForbidden)
      return
    }
  }
}

//updateProfileByAdmin is HandlerFunc that can be mapped to requests for deleting users. This function can
//be used for deleting admins and users. Should be performed only for requests from users with admin role.
func deleteProfileByAdmin(service controllers.UserServiceInterface) http.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request) {
    idStr := mux.Vars(r)["user_id"]
    id, err := strconv.Atoi(idStr)

    if err != nil {
      responseError(w, err, http.StatusForbidden)
      return
    }

    respUserDto, err := service.DeleteUser(int64(id))
    if err != nil {
      responseError(w, err, http.StatusForbidden)
      return
    }

    respBody, err := json.Marshal(respUserDto)
    _, err = w.Write(respBody)
    if err != nil {
      responseError(w, err, http.StatusForbidden)
      return
    }
  }
}
