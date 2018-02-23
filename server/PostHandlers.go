package server

import (
  "net/http"
  "encoding/json"
  "fmt"
  "github.com/gorilla/mux"
  "strconv"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
  "github.com/aws/aws-sdk-go/service/s3/s3iface"
  "go-team-room/controllers"
  "time"
)

//To CREATE new post in DynamoDB Table "Post"
func CreateNewPost(svcd dynamodbiface.DynamoDBAPI, svcs s3iface.S3API) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {

    var post controllers.Post

    //Decode request MULTIPART/FORM-DATA
    r.ParseMultipartForm(0)
    post.Title = r.FormValue("post_title")
    post.Text = r.FormValue("post_text")
    post.UserID = r.FormValue("user_id")

    //Set "post_id", "post_like", "file_link"
    post.PostID = time.Now().String()
    post.LastUpdate = post.PostID
    post.Like = make([]*string,0)
    likeInit := "NULL"
    post.Like = append(post.Like, &likeInit)
    post.FileLink = "NULL"
    //Check if file exists in request
    //if exists UPLOAD to S3
    //if not - "file_link" remains "NULL"
    if fhs := r.MultipartForm.File["upfile"]; len(fhs) > 0 {

      //Get File, File Header, Error from Multipart Form File
      file, handler, err := r.FormFile("upfile")

      if err != nil {
        log.Debug("Error to get file")
        return
      }

      defer file.Close()

      //UPLOAD file to S3 and GET link
      post.FileLink = controllers.UploadFileToS3(svcs, file, handler)
    }

    resp, err := controllers.CreateNewPost(svcd, post)

    if err != nil {
      w.WriteHeader(http.StatusNoContent)
    }
    log.Debug("Post created")
    //Encode response JSON
    _ = json.NewEncoder(w).Encode(&resp)
  }
}

//To DELETE existing post by "post_id" from DynamoDB Table "Post"
func DeletePost(svcd dynamodbiface.DynamoDBAPI, svcs s3iface.S3API) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)

    response := controllers.DeletePost(svcd, svcs, vars["post_id"])

    //Encode response JSON
    _ = json.NewEncoder(w).Encode(&response)

  }
}

//To GET post by "post_id" from DynamoDB Table "Post"
func GetPost(svc dynamodbiface.DynamoDBAPI) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    //Gorilla tool to handle request "/post/{post_id}" with method GET
    vars := mux.Vars(r)

    post, err := controllers.GetPost(svc, vars["post_id"])

    if err != nil {
      w.WriteHeader(http.StatusNoContent)
      return
    }

    //Encode response JSON
    _ = json.NewEncoder(w).Encode(&post)
  }
}

//To GET posts by "user_id" from DynamoDB Table "Post"
func GetPostByUserID(svc dynamodbiface.DynamoDBAPI) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    post, err := controllers.GetPostByUserID(svc, vars["user_id"])
    log.Info(post)
    if err != nil {
      w.WriteHeader(http.StatusNoContent)
    }

    _ = json.NewEncoder(w).Encode(&post)
  }
}

//To UPDATE post in DynamoDB Table "Post"
func UpdatePost (svc dynamodbiface.DynamoDBAPI) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    var post controllers.Post

    //Decode request JSON
    _ = json.NewDecoder(r.Body).Decode(&post)
    post.LastUpdate = time.Now().String()

    //Gorilla tool to handle "/post/{post_id}" with method PUT
    vars := mux.Vars(r)
    post.PostID = vars["post_id"]

    post, _ = controllers.UpdatePost(svc, post)

    //Encode response JSON
    _ = json.NewEncoder(w).Encode(&post)
  }
}

//To GET file from S3
func GetFileFromS3(svc s3iface.S3API) http.HandlerFunc  {
  return func(w http.ResponseWriter, r *http.Request) {

    //Gorilla tool to handle request "/post/{post_id}" with method DELETE
    vars := mux.Vars(r)
    fileName := vars["file_link"]
    buff, err := controllers.DownloadFileFromS3(svc, fileName)

    if err != nil {
      fmt.Println("Error")
    }

    //Get the Content-Type of the file
    //Create a buffer to store the header of the file in
    FileHeader := make([]byte, 512)
    //Copy the headers into the FileHeader buffer
    FileHeader = buff.Bytes()[:512]
    //Get content type of file
    FileContentType := http.DetectContentType(FileHeader)

    FileSize := strconv.FormatInt(int64(len(buff.Bytes())), 10) //Get file size as a string

    //Send the headers
    w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
    w.Header().Set("Content-Type", FileContentType)
    w.Header().Set("Content-Length", FileSize)

    w.Write(buff.Bytes())
  }
}

func SetLike(svc dynamodbiface.DynamoDBAPI) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    log.Println(r)

  }
}
