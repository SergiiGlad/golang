package server

import (
  "net/http"
  "io/ioutil"
  "log"
  "fmt"
  "encoding/json"
  "go-team-room/models/dto"
  "go-team-room/models"
  "go-team-room/controllers"
  "github.com/gorilla/mux"
  "strconv"
)

func createProfile(w http.ResponseWriter, r *http.Request) {

  body, err := ioutil.ReadAll(r.Body)

  if err != nil {
    responseError(w, err)
    return
  }

  user := dto.RequestUserDto{}
  err = json.Unmarshal(body, &user)

  if err != nil {
    log.Println(err)
    responseError(w, err)
    return
  }

  userDao := models.RequestUserDtoToDao(user)

  err = controllers.CreateUser(&userDao)

  if err != nil {
    log.Println(err)
    responseError(w, err)
    return
  }

  respUserDto := models.UserDaoToResponseDto(userDao)

  body, err = json.Marshal(respUserDto)
  _, err = w.Write(body)

  if err != nil {
    log.Println(err)
    responseError(w, err)
    return
  }
}

func updateProfile(w http.ResponseWriter, r *http.Request) {

  body, err := ioutil.ReadAll(r.Body)

  if err != nil {
    responseError(w, err)
    return
  }

  user := dto.RequestUserDto{}
  err = json.Unmarshal(body, &user)

  if err != nil {
    log.Println(err)
    responseError(w, err)
    return
  }

  userDao := models.RequestUserDtoToDao(user)
  idStr := mux.Vars(r)["id"]
  id, err := strconv.Atoi(idStr)

  if err != nil {
    log.Println(err)
    responseError(w, err)
    return
  }

  err = controllers.UpdateUser(int64(id), &userDao)

  if err != nil {
    log.Println(err)
    responseError(w, err)
    return
  }

  respUserDto := models.UserDaoToResponseDto(userDao)

  body, err = json.Marshal(respUserDto)
  _, err = w.Write(body)

  if err != nil {
    log.Println(err)
    responseError(w, err)
    return
  }
}

func deleteProfile(w http.ResponseWriter, r *http.Request) {
  idStr := mux.Vars(r)["id"]
  id, err := strconv.Atoi(idStr)

  if err != nil {
    log.Println(err)
    responseError(w, err)
    return
  }

  userDao, err := controllers.DeleteUser(int64(id))

  if err != nil {
    log.Println(err)
    responseError(w, err)
    return
  }

  respUserDto := models.UserDaoToResponseDto(*userDao)

  body, err := json.Marshal(respUserDto)
  _, err = w.Write(body)

  if err != nil {
    log.Println(err)
    responseError(w, err)
    return
  }
}

func responseError(w http.ResponseWriter, err error) {
  rerror := dto.ResponseError{err.Error()}

  body, err := json.Marshal(rerror)
  if err != nil {
    log.Println(err)
    fmt.Fprint(w, err)
  }

  w.WriteHeader(http.StatusBadRequest)
  w.Write(body)
}
