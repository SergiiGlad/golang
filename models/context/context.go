package context

import (
  "github.com/gorilla/context"
  "net/http"
)

type contextKey int

const (
  user_id contextKey = 0
  email contextKey = 1
  firstName contextKey = 2
  lastName contextKey = 3
  phone contextKey = 4
  userRole contextKey = 5
  // if need to add NewKey for context just add it here
  // and implement GetNewKeyFromContext(r *http.Request) SomeType{...} function
  // and SetNewKeyToContext(r *http.Request, val SomeType) {...} function
)

func GetIdFromContext(r *http.Request) int64 {
  if val := context.Get(r, user_id); val != nil {
    return val.(int64)
  }
  return 0
}

func GetEmailFromContext(r *http.Request) string {
  if val := context.Get(r, email); val != nil {
    return val.(string)
  }
  return ""
}

func GetFirstNameFromContext(r *http.Request) string {
  if val := context.Get(r, firstName); val != nil {
    return val.(string)
  }
  return ""
}

func GetLastNameFromContext(r *http.Request) string {
  if val := context.Get(r, lastName); val != nil {
    return val.(string)
  }
  return ""
}

func GetPhoneFromContext(r *http.Request) string {
  if val := context.Get(r, phone); val != nil {
    return val.(string)
  }
  return ""
}

func GetRoleFromContext(r *http.Request) string {
  if val := context.Get(r, userRole); val != nil {
    return val.(string)
  }
  return ""
}

func SetIdToContext(r *http.Request, val int64) {
  context.Set(r, user_id, val)
}

func SetEmailToContext(r *http.Request, val string) {
  context.Set(r, email, val)
}

func SetFirstNameToContext(r *http.Request, val string) {
  context.Set(r, firstName, val)
}

func SetLastNameToContext(r *http.Request, val string) {
  context.Set(r, lastName, val)
}

func SetPhoneToContext(r *http.Request, val string) {
  context.Set(r, phone, val)
}

func SetUserRoleToContext(r *http.Request, val string) {
  context.Set(r, userRole, val)
}
