package server

import (
  "net/http"
  "fmt"
)

func logout(w http.ResponseWriter, r *http.Request) {
  session, err := store.Get(r, "name")

  if err != nil {
    responseError(w, err, http.StatusBadRequest)
    return
  }

  if session.Values["loginned"] == false {
    return
  }

  session.Values["loginned"] = false
  session.Options.MaxAge = -1
  session.Save(r, w)
  fmt.Fprintf(w, "User was logout")
}
