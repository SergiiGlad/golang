package server

import (
  "go-team-room/models/dto"
  "io/ioutil"
  "net/http"
  "encoding/json"
)

//responseError is helper function that can be used for returning error for client. It accepts
// http.ResponseWriter where response should be written. err is error that response body should contain.
// error message will be written as a value for "reason" key in json type. Function also needs response
// http code.
func responseError(w http.ResponseWriter, err error, code int) {
  rerror := dto.ResponseError{err.Error()}
  log.Println(err)

  body, err := json.Marshal(rerror)
  if err != nil {
    log.Println(err)
  }

  http.Error(w, string(body), code)
}

//userDtoFromReq read *http.Request body and tries to unmarshal body content into dto.RequestUserDto.
//If unmarshal operation performed successfully error will be nil (don't forget to check it).
func userDtoFromReq(request *http.Request) (dto.RequestUserDto, error) {
  body, err := ioutil.ReadAll(request.Body)
  userDto := dto.RequestUserDto{}
  if err != nil {
    return userDto, err
  }

  err = json.Unmarshal(body, &userDto)
  if err != nil {
    return userDto, err
  }

  return userDto, err
}

//dtoFromReq read *http.Request body and tries to unmarshal body content into golang type you passed
//as dto argument. If unmarshal operation performed successfully error will be nil (don't forget to check it).
//After function successfully performed you need to make type assertion.
func dtoFromReq(request *http.Request, dto interface{}) error {
  body, err := ioutil.ReadAll(request.Body)
  if err != nil {
    return err
  }

  err = json.Unmarshal(body, &dto)
  if err != nil {
    return err
  }

  return nil
}
