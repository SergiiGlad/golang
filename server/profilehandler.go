package server

import (
  "go-team-room/controllers"
  "github.com/aws/aws-sdk-go/service/s3/s3iface"
  "net/http"
  "go-team-room/models/dto"
  "github.com/gorilla/mux"
  "strconv"
  "encoding/json"
  "io"
  "strings"
)

func GetProfile(service controllers.UserServiceInterface) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    idStr := mux.Vars(r)["user_id"]
    id, err := strconv.Atoi(idStr)

    if err != nil {
      log.Error(err)
      return
    }

    responseUserDTO, err := service.GetUser(int64(id))

    if err != nil {
      log.Error(err)
      return
    }

    _ = json.NewEncoder(w).Encode(&responseUserDTO)
  }
}

func UploadAvatar (service controllers.UserServiceInterface, svc s3iface.S3API) http.HandlerFunc{
  return func(w http.ResponseWriter, r *http.Request) {

    var userDTO dto.RequestUserDto

    r.ParseMultipartForm(0)
    if fhs := r.MultipartForm.File["upfile"]; len(fhs) > 0 {

      //Get File, File Header, Error from Multipart Form File
      file, handler, err := r.FormFile("upfile")
      buff := make([]byte, 512)
      n, err := file.Read(buff)

      if err != nil && err != io.EOF {
        log.Debug("Error")
        return
      }

      file.Seek(0,0)

      contentType := http.DetectContentType(buff[:n])

      if !strings.HasPrefix(contentType, "image") {
        w.WriteHeader(http.StatusBadRequest)
        return
      }

      defer file.Close()

      //UPLOAD file to S3 and GET link
      userDTO.Avatar = controllers.UploadFileToS3(svc, file, handler)
    }

    idStr := mux.Vars(r)["user_id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
      responseError(w, err, http.StatusForbidden)
      return
    }

    respUserDto, err := service.UpdateUser(int64(id), &userDTO)
    if err != nil {
      responseError(w, err, http.StatusForbidden)
      return
    }

    _ = json.NewEncoder(w).Encode(&respUserDto)

  }
}

func DeleteAvatar (service controllers.UserServiceInterface, svc s3iface.S3API) http.HandlerFunc{
  return func(w http.ResponseWriter, r *http.Request) {

    idStr := mux.Vars(r)["user_id"]
    id, err := strconv.Atoi(idStr)

    if err != nil {
      log.Error(err)
      return
    }

    userDTO, err :=  service.GetUser(int64(id))


    controllers.DeleteFileFromS3(userDTO.Avatar, svc)

    var requestUserDTO dto.RequestUserDto

    requestUserDTO.Avatar = "NULL"

    responseUserDTO, err := service.UpdateUser(int64(id), &requestUserDTO)

    if err != nil {
      log.Error(err)
      return
    }

    _ = json.NewEncoder(w).Encode(&responseUserDTO)
  }
}
