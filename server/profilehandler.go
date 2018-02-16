package server

import (
  "go-team-room/controllers"
  "github.com/aws/aws-sdk-go/service/s3/s3iface"
  "net/http"
  "go-team-room/models/dto"
  "github.com/gorilla/mux"
  "strconv"
  "encoding/json"
)

func UploadAvatar (service controllers.UserServiceInterface, svc s3iface.S3API) http.HandlerFunc{
  return func(w http.ResponseWriter, r *http.Request) {

    var userDTO dto.RequestUserDto

    r.ParseMultipartForm(0)

    if fhs := r.MultipartForm.File["upfile"]; len(fhs) > 0 {

      //Get File, File Header, Error from Multipart Form File
      file, handler, err := r.FormFile("upfile")

      if err != nil {
        log.Debug("Error to get file")
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
