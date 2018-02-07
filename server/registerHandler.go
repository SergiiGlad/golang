package server

import (
  "net/http"
  "encoding/json"
  "go-team-room/controllers"
  "go-team-room/models/dto"
  "go-team-room/models/dao/entity"
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
    log.Println(respUserDto)
  }
}