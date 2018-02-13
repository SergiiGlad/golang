package server

import (
  "net/http"
  "github.com/gorilla/mux"
  "go-team-room/controllers"
  "github.com/pkg/errors"
)

func ConfirmAccount(service controllers.TokenGeneratorInterface) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    token := vars["token"]
    approved, err := service.ApproveUser(token)
    if err != nil {
      log.Errorf("Failed to approve user for token: %s, error: %s", token, err)
      responseError(w, err, http.StatusBadRequest)
      return
    }
    if !approved {
      log.Warnf("Cant approve user for token: %s", token)
      responseError(w, errors.New("Cant approve user for token this token." ), http.StatusBadRequest)
      return
    }
    log.Infof("Successfully approve user for token: %s", token)
  }
}
