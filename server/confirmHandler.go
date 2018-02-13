package server

import (
  "net/http"
  "github.com/gorilla/mux"
  "go-team-room/controllers"
)

func ConfirmAccount(service controllers.TokenGeneratorInterface) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    token := vars["token"]
    approved, err := service.ApproveUser(token)
    if err != nil {
      log.Error("Failed to approve user for token: {}, error", token)
      responseError(w, err, http.StatusBadRequest)
      return
    }
    if !approved {
      log.Warn("Cant approve user for token: {}", token)
      responseError(w, err, http.StatusBadRequest)
      return
    }
    log.Info("Successfully approve user for token: {}", token)
  }
}
