package server

import (
  "github.com/gorilla/sessions"
  "net/http"
)

var (
  store = sessions.NewCookieStore([]byte("wesomekey"))
)

func Authorize(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
      if r.URL.Path == "/login" || r.URL.Path == "/registration" || r.URL.Path == "/" || r.URL.Path == "/logout" ||
        r.URL.Path == "/confirm/email/{token}" {
        next.ServeHTTP(w, r)
      } else {
        http.Redirect(w, r, "/login", 301)
      }
    }
  })
}
