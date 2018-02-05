package server

import (
  "net/http"
  "encoding/json"
  "go-team-room/controllers"
  //"go-team-room/models/dao"
  "go-team-room/models/dto"
  "go-team-room/models/dao/entity"
  //"github.com/gorilla/mux"
  //"io/ioutil"
  //"fmt"
  //"strings"
)

func ProtectionUserRole(userDto *dto.RequestUserDto) {
  if userDto.Role == entity.AdminRole{
    userDto.Role = entity.UserRole
  }
}

func registerUser(service controllers.UserServiceInterface) http.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request) {
    reqUserDto, err := userDtoFromReq(r)

    ProtectionUserRole(&reqUserDto)

    if err != nil {
      responseError(w, err, http.StatusForbidden)
      return
    }

    respUserDto, err := service.CreateUser(&reqUserDto)

    if err != nil {
      responseError(w, err, http.StatusForbidden)
      return
    }

    respBody, err := json.Marshal(respUserDto)
    _, err = w.Write(respBody)

    if err != nil {
      responseError(w, err, http.StatusForbidden)
      return
    }
  }
}


//func recoveryPass(w http.ResponseWriter, r *http.Request) {
//    r.ParseForm()       // parse arguments, you have to call this by yourself
//    fmt.Println(r.Form) // print form information in server side
//    fmt.Println("path", r.URL.Path)
//    fmt.Println("scheme", r.URL.Scheme)
//    fmt.Println(r.Form["url_long"])
//    for k, v := range r.Form {
//      fmt.Println("key:", k)
//      fmt.Println("val:", strings.Join(v, ""))
//    }
//    fmt.Fprintf(w, "Hello astaxie!") // send data to client side
//}
