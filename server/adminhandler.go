package server

import (
  "net/http"
  "io/ioutil"
  "encoding/json"
  "go-team-room/models/dto"
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

//responseError is helper function that can be used for returning error for client. It accepts
// http.ResponseWriter where response should be written. err is error that response body should contain.
// error message will be written as a value for "reason" key in json type. Function also needs response
// http code.
func responseError(w http.ResponseWriter, err error, code int) {
  rerror := dto.ResponseError{err.Error()}
  log.Println(err)

  body, err := json.Marshal(rerror)
  if err != nil {
    log.Println(err)
  }

  http.Error(w, string(body), code)
}

//userDtoFromReq read *http.Request body and tries to unmarshal body content into dto.RequestUserDto.
//If unmarshal operation performed successfully error will be nil (don't forget to check it).
func userDtoFromReq(request *http.Request) (dto.RequestUserDto, error) {
  body, err := ioutil.ReadAll(request.Body)
  userDto := dto.RequestUserDto{}

  if err != nil {
    return userDto, err
  }

  err = json.Unmarshal(body, &userDto)

  if err != nil {
    return userDto, err
  }

  return userDto, err
}
