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
      session.Options.MaxAge = -1
      session.Save(r, w)
      responseError(w, err, http.StatusForbidden)
      return
    }

    if session.Values["loginned"] == true {
      next.ServeHTTP(w, r)
    } else {
      // add session flashes when UI will be ready
      paths := []string{"/", "/login", "/logout", "/registration"} //paths no needed to check authorization

      for _, path := range paths {
        if r.URL.Path == path {
          next.ServeHTTP(w, r)
          return
        }
      }

      http.Error(w, "Forbidden", http.StatusForbidden)
    }
  })
}
