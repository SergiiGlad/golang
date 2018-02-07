package server

import (
  "go-team-room/controllers"
  "net/http"
  "go-team-room/models/dto"
  "github.com/gorilla/mux"
  "strconv"
)

func getFriends(service controllers.FriendServiceInterface) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    var friendDto dto.Friendship
    err := dtoFromReq(r, &friendDto)

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

    service.GetUserFriends(int64(id))

  }
}
