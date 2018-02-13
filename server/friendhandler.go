package server

import (
  "go-team-room/controllers"
  "net/http"
  "github.com/gorilla/mux"
  "strconv"
  "encoding/json"
  "go-team-room/models/dao/entity"
  "github.com/pkg/errors"
)

func getFriends(service controllers.FriendServiceInterface) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    idStr := mux.Vars(r)["user_id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
      responseError(w, err, http.StatusBadRequest)
      return
    }

    friends, err := service.GetFriends(int64(id))
    if err != nil {
      responseError(w, err, http.StatusBadRequest)
    }

    respBody, err := json.Marshal(friends)
    _, err = w.Write(respBody)
    if err != nil {
      responseError(w, err, http.StatusBadRequest)
      return
    }
  }
}

func getUsersWithRequests(service controllers.FriendServiceInterface) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    idStr := mux.Vars(r)["user_id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
      responseError(w, err, http.StatusBadRequest)
      return
    }

    friends, err := service.GetUsersWithRequests(int64(id))
    if err != nil {
      responseError(w, err, http.StatusBadRequest)
    }

    respBody, err := json.Marshal(friends)
    _, err = w.Write(respBody)
    if err != nil {
      responseError(w, err, http.StatusBadRequest)
      return
    }
  }
}

func newFriendRequest(service controllers.FriendServiceInterface) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    var connection entity.Connection
    err := dtoFromReq(r, &connection)
    if err != nil {
      responseError(w, err, http.StatusBadRequest)
      return
    }

    err = service.NewFriendRequest(&connection)
    if err != nil {
      responseError(w, err, http.StatusBadRequest)
      return
    }

    w.WriteHeader(http.StatusOK)
  }
}

func replyToFriendRequest(service controllers.FriendServiceInterface) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    var connection entity.Connection
    err := dtoFromReq(r, &connection)
    if err != nil {
      responseError(w, err, http.StatusBadRequest)
      return
    }

    if connection.Status == entity.Approved {
      err = service.ApproveFriendRequest(&connection)
    } else if connection.Status == entity.Rejected {
      err = service.RejectFriendRequest(&connection)
    } else {
      responseError(w, errors.New("Invalid request reply"), http.StatusBadRequest)
    }
    if err != nil {
      responseError(w, err, http.StatusBadRequest)
      return
    }

    w.WriteHeader(http.StatusOK)
  }
}

func deleteFriendship(service controllers.FriendServiceInterface) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    var connection entity.Connection
    err := dtoFromReq(r, &connection)
    if err != nil {
      responseError(w, err, http.StatusBadRequest)
      return
    }

    err = service.DeleteFriendship(&connection)
    if err != nil {
      responseError(w, err, http.StatusBadRequest)
      return
    }

    w.WriteHeader(http.StatusOK)
  }
}
