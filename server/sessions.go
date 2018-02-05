package server

import (
  "github.com/gorilla/sessions"
  "net/http"
)

var (
  store = sessions.NewCookieStore([]byte("secretkey"))
)

func Authorize(next http.Handler) http.Handler {
  return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
    session, err := store.Get(r, "name")

    if err != nil {
      responseError(w, err)
      return
    }
    if session.Values["auth"] == "loginned" {
      next.ServeHTTP(w, r)
    } else {
      if r.URL.Path == "/login" {
        next.ServeHTTP(w, r)
      } else {
        http.Redirect(w, r, "/login", 301)
      }
    }
  })
}
