package server

import (
  "net/http"
  "encoding/json"
  "fmt"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/aws/aws-sdk-go/aws/awserr"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
  "github.com/gorilla/mux"
  "github.com/aws/aws-sdk-go/service/dynamodb/expression"
  "os"
  "time"
  "mime/multipart"
  "github.com/aws/aws-sdk-go/service/s3/s3manager"
  "io"
  "crypto/rand"
  "strings"
  "github.com/aws/aws-sdk-go/service/s3"
  "strconv"
)

//Post structure
type Post struct{
  Title string `json:"post_title"`
  Text string	`json:"post_text"`
  PostID string `json:"post_id"`
  UserID string `json:"user_id"`
  Like string `json:"post_like"`
  FileLink string `json:"file_link"`
}

//To CREATE new post in DynamoDB Table "Post"
func CreateNewPost(w http.ResponseWriter, r *http.Request) {

  var newPost Post

  //Create new Session for DynamoDB
  sess, err := session.NewSession(&aws.Config{
    Region: aws.String("eu-west-2"),
  })

  //Decode request MULTIPART/FORM-DATA
  r.ParseMultipartForm(0)
  newPost.Title = r.FormValue("post_title")
  newPost.Text = r.FormValue("post_text")
  newPost.UserID = r.FormValue("user_id")

  //Set "post_id", "post_like", "file_link"
  newPost.PostID = time.Now().String()
  newPost.Like = "0"
  newPost.FileLink = "NULL"

  //Check if file exists in request
  //if exists UPLOAD to S3
  //if not - "file_link" remains "NULL"
  if fhs := r.MultipartForm.File["upfile"]; len(fhs) > 0{

    //Get File, File Header, Error from Multipart Form File
    file, handler, err := r.FormFile("upfile")

    if err != nil {
      fmt.Println(err)
      return
    }

    defer file.Close()

    //UPLOAD file to S3 and GET link
    newPost.FileLink = UploadFileToS3(sess, file, handler)
  }

  svc := dynamodb.New(sess)

  //Request to DynamoDB to CREATE new post with KEY_ATTRIBUTE "post_id" (TimeStamp)
  input := &dynamodb.PutItemInput{
    Item: map[string]*dynamodb.AttributeValue{
      "post_id": {
        S: &newPost.PostID,
      },
      "post_title": {
        S: &newPost.Title,
      },
      "post_text": {
        S: &newPost.Text,
      },
      "user_id": {
        S: &newPost.UserID,
      },
      "post_like": {
        N: &newPost.Like,
      },
      "file_link": {
        S: &newPost.FileLink,
      },
    },
    ReturnConsumedCapacity: aws.String("TOTAL"),
    TableName:              aws.String("Post"), //Name of Table in DynamoDB
  }

  //Get result
  result, err := svc.PutItem(input)

  //ERROR block (from AWS SDK GO DynamoDB documentation)
  if err != nil {
    if aerr, ok := err.(awserr.Error); ok {
      switch aerr.Code() {
      case dynamodb.ErrCodeConditionalCheckFailedException:
        fmt.Println(dynamodb.ErrCodeConditionalCheckFailedException, aerr.Error())
      case dynamodb.ErrCodeProvisionedThroughputExceededException:
        fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
      case dynamodb.ErrCodeResourceNotFoundException:
        fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
      case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
        fmt.Println(dynamodb.ErrCodeItemCollectionSizeLimitExceededException, aerr.Error())
      case dynamodb.ErrCodeInternalServerError:
        fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
      default:
        fmt.Println(aerr.Error())
      }
    } else {
      // Print the error, cast err to awserr.Error to get the Code and
      // Message from an error.
      fmt.Println(err.Error())
    }
    return
  }

  //Encode response JSON
  _ = json.NewEncoder(w).Encode(&newPost)

  //Print response in console
  fmt.Println(result)
}

//To DELETE existing post by "post_id" from DynamoDB Table "Post"
func DeletePost(w http.ResponseWriter, r *http.Request) {
  var post Post

  //Gorilla tool to handle request "/post/{post_id}" with method DELETE
  vars := mux.Vars(r)
  post.PostID = vars["post_id"]

  //Create new Session for DynamoDB
  sess, err := session.NewSession(&aws.Config{
    Region:      aws.String("eu-west-2"),
  })

  post = GetPostBody(post, sess)

  if post.FileLink != "NULL"{
    DeleteFileFromS3(post.FileLink, sess)
  }
  svc := dynamodb.New(sess)

  //Request to DynamoDB to DELETE post with KEY_ATTRIBUTE "post_id" (TimeStamp)
  input := &dynamodb.DeleteItemInput{
    Key: map[string]*dynamodb.AttributeValue{
      "post_id": {
        S: &post.PostID,
      },
    },
    TableName: aws.String("Post"),
  }

  //Get result
  result, err := svc.DeleteItem(input)

  //ERROR block (from AWS SDK GO DynamoDB documentation)
  if err != nil {
    if aerr, ok := err.(awserr.Error); ok {
      switch aerr.Code() {
      case dynamodb.ErrCodeConditionalCheckFailedException:
        fmt.Println(dynamodb.ErrCodeConditionalCheckFailedException, aerr.Error())
      case dynamodb.ErrCodeProvisionedThroughputExceededException:
        fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
      case dynamodb.ErrCodeResourceNotFoundException:
        fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
      case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
        fmt.Println(dynamodb.ErrCodeItemCollectionSizeLimitExceededException, aerr.Error())
      case dynamodb.ErrCodeInternalServerError:
        fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
      default:
        fmt.Println(aerr.Error())
      }
    } else {
      // Print the error, cast err to awserr.Error to get the Code and
      // Message from an error.
      fmt.Println(err.Error())
    }
    return
  }

  //Encode response JSON
  _ = json.NewEncoder(w).Encode(&post.PostID)

  //Print response in console
  fmt.Println(result)
}

//To GET post by "post_id" from DynamoDB Table "Post"
func GetPost(w http.ResponseWriter, r *http.Request) {
  var post Post

  //Gorilla tool to handle request "/post/{post_id}" with method GET
  vars := mux.Vars(r)
  post.PostID = vars["post_id"]

  //Create new session for DynamoDB
  sess, err := session.NewSession(&aws.Config{
    Region:      aws.String("eu-west-2"),
  })
  svc := dynamodb.New(sess)

  //Request to DynamoDB to GET post by "post_id"
  input := &dynamodb.GetItemInput{
    Key: map[string]*dynamodb.AttributeValue{
      "post_id": {
        S: &post.PostID,
      },
    },
    TableName: aws.String("Post"),
  }

  //Get result
  result, err := svc.GetItem(input)

  //ERROR block (from AWS SDK GO DynamoDB documentation)
  if err != nil {
    if aerr, ok := err.(awserr.Error); ok {
      switch aerr.Code() {
      case dynamodb.ErrCodeProvisionedThroughputExceededException:
        fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
      case dynamodb.ErrCodeResourceNotFoundException:
        fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
      case dynamodb.ErrCodeInternalServerError:
        fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
      default:
        fmt.Println(aerr.Error())
      }
    } else {
      // Print the error, cast err to awserr.Error to get the Code and
      // Message from an error.
      fmt.Println(err.Error())
    }
    return
  }

  //Check if POST exists in table, if not: RESPONSE 204 - "No Content"
  if len(result.Item) == 0 {
    w.WriteHeader(http.StatusNoContent)
    return
  }

  //Unmarshal result to Post structure
  err = dynamodbattribute.UnmarshalMap(result.Item, &post)

  //Print result in console
  fmt.Println(result)

  //Encode response JSON
  _ = json.NewEncoder(w).Encode(&post)
}

//To GET posts by "user_id" from DynamoDB Table "Post"
func GetPostByUserID(w http.ResponseWriter, r *http.Request){
  var post Post

  //Gorilla tool to handle request "/post/{user_id}" with method GET
  vars := mux.Vars(r)
  post.UserID = vars["user_id"]

  //Create new session for DynamoDB
  sess, err := session.NewSession(&aws.Config{
    Region:      aws.String("eu-west-2"),
  })
  svc := dynamodb.New(sess)

  //Filter expression: Seeks all items in table with equal "user_id"
  filt := expression.Name("user_id").Equal(expression.Value(post.UserID))

  //Make projection: displays all expression.Name with equal "user_id"
  proj := expression.NamesList(expression.Name("post_title"), expression.Name("post_text"), expression.Name("post_id"), expression.Name("user_id"), expression.Name("post_like"), expression.Name("file_link"))

  //Build expression with filter and projection
  expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

  //Parameters for expression
  params := &dynamodb.ScanInput{
    ExpressionAttributeNames:  expr.Names(),
    ExpressionAttributeValues: expr.Values(),
    FilterExpression:          expr.Filter(),
    ProjectionExpression:      expr.Projection(),
    TableName:                 aws.String("Post"),
  }

  //Get result
  result, err := svc.Scan(params)

  //Check if POST exists in table, if not: RESPONSE 204 - "No Content"
  if len(result.Items) == 0 {
    w.WriteHeader(http.StatusNoContent)
    return
  }

  //Loop for encoding multiples post which was made by "user_id"
  for _, i := range result.Items {
    post := Post{}

    //Unmarshal result to Post structure
    err = dynamodbattribute.UnmarshalMap(i, &post)

    //Check error for unmarshalling
    if err != nil {
      fmt.Println("Got error unmarshalling:")
      fmt.Println(err.Error())
      os.Exit(1)
    }

    //Encode response JSON
    _ = json.NewEncoder(w).Encode(&post)
  }

  //Print result in console
  fmt.Print(result)
}

//To UPDATE post in DynamoDB Table "Post"
func UpdatePost(w http.ResponseWriter, r *http.Request){
  var post Post

  //Decode request JSON
  _ = json.NewDecoder(r.Body).Decode(&post)

  //Gorilla tool to handle "/post/{post_id}" with method PUT
  vars := mux.Vars(r)
  post.PostID = vars["post_id"]

  //Create new session for DynamoDB
  sess, err := session.NewSession(&aws.Config{
    Region: aws.String("eu-west-2"),
  })
  svc := dynamodb.New(sess)

  //Request to DynamoDB to UPDATE Item in table
  input := &dynamodb.UpdateItemInput{
    ExpressionAttributeNames: map[string]*string{
      "#PTitle": aws.String("post_title"),
      "#PText": aws.String("post_text"),
    },
    //Attributes to UPDATE
    ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
      ":t": {
        S: &post.Title,
      },
      ":e": {
        S: &post.Text,
      },
    },
    Key: map[string]*dynamodb.AttributeValue{
      "post_id": {
        S: &post.PostID,
      },
    },
    ReturnValues:     aws.String("UPDATED_NEW"),
    TableName: aws.String("Post"),
    UpdateExpression: aws.String("set #PTitle = :t, #PText = :e"),
  }

  //Get result
  result, err := svc.UpdateItem(input)

  //ERROR block (from AWS SDK GO DynamoDB documentation)
  if err != nil {
    if aerr, ok := err.(awserr.Error); ok {
      switch aerr.Code() {
      case dynamodb.ErrCodeConditionalCheckFailedException:
        fmt.Println(dynamodb.ErrCodeConditionalCheckFailedException, aerr.Error())
      case dynamodb.ErrCodeProvisionedThroughputExceededException:
        fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
      case dynamodb.ErrCodeResourceNotFoundException:
        fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
      case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
        fmt.Println(dynamodb.ErrCodeItemCollectionSizeLimitExceededException, aerr.Error())
      case dynamodb.ErrCodeInternalServerError:
        fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
      default:
        fmt.Println(aerr.Error())
      }
    } else {
      // Print the error, cast err to awserr.Error to get the Code and
      // Message from an error.
      fmt.Println(err.Error())
    }
    return
  }

  //Print result in console
  fmt.Println(result)

  //Encode response JSON
  _ = json.NewEncoder(w).Encode(&post)

}

//To UPLOAD file to S3
func UploadFileToS3 (sess *session.Session, f multipart.File, handl *multipart.FileHeader) string{

  // Create an uploader with the session and default options
  uploader := s3manager.NewUploader(sess)

  fileType := handl.Filename[strings.LastIndexAny(handl.Filename, "."):]
  fileType = strings.ToLower(fileType)

  //Generate UUID for File
  uuid, err := newUUID()
  if err != nil {
    fmt.Printf("error: %v\n", err)
  }

  // Upload the file to S3.
  result, err := uploader.Upload(&s3manager.UploadInput{
    Bucket: aws.String("gohumfiles"),
    Key:    aws.String(uuid + fileType),
    Body:   f,
  })
  if err != nil {
    fmt.Errorf("failed to upload file, %v", err)
    return "failed to upload file"
  }
  fmt.Printf("file uploaded to, %s\n", aws.StringValue(&result.Location))
  link := "/uploads/" + uuid + fileType
  return link
}

//To GET file from S3
func GetFileFromS3(w http.ResponseWriter, r *http.Request) {

  //Gorilla tool to handle request "/post/{post_id}" with method DELETE
  vars := mux.Vars(r)
  fileName := vars["file_link"]

  //Create new Session for DynamoDB
  sess, err := session.NewSession(&aws.Config{
    Region: aws.String("eu-west-2"),
  })

  // Create a downloader with the session and default options
  downloader := s3manager.NewDownloader(sess)

  // Create a file to write the S3 Object contents to.
  f, err := os.Create(fileName)
  if err != nil {
    fmt.Errorf("failed to create file %q, %v", fileName, err)
    return
  }

  // Write the contents of S3 Object to the file
  n, err := downloader.Download(f, &s3.GetObjectInput{
    Bucket: aws.String("gohumfiles"),
    Key:    aws.String(fileName),
  })
  if err != nil {
    fmt.Errorf("failed to download file, %v", err)
    return
  }
  fmt.Printf("file downloaded, %d bytes\n", n)

  f, err = os.Open(fileName)
  defer f.Close()

  if err != nil {
    //File not found, send 404
    http.Error(w, "File not found.", 404)
    return
  }

  //Get the Content-Type of the file
  //Create a buffer to store the header of the file in
  FileHeader := make([]byte, 512)
  //Copy the headers into the FileHeader buffer
  f.Read(FileHeader)
  //Get content type of file
  FileContentType := http.DetectContentType(FileHeader)

  //Get the file size
  FileStat, _ := f.Stat()                     //Get info from file
  FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string

  //Send the headers
  w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
  w.Header().Set("Content-Type", FileContentType)
  w.Header().Set("Content-Length", FileSize)

  //Send the file
  //We read 512 bytes from the file already so we reset the offset back to 0
  f.Seek(0, 0)
  io.Copy(w, f) //'Copy' the file to the client

  os.Remove(fileName)
}

//To DELETE file from S3 when DELETE Post
func DeleteFileFromS3(file string, sess *session.Session) {
  filename := strings.TrimPrefix(file, "/uploads/")

  svc := s3.New(sess)

  input := &s3.DeleteObjectsInput{
    Bucket: aws.String("gohumfiles"),
    Delete: &s3.Delete{
      Objects: []*s3.ObjectIdentifier{
        {
          Key: aws.String(filename),
        },
      },
      Quiet: aws.Bool(false),
    },
  }

  result, err := svc.DeleteObjects(input)
  if err != nil {
    if aerr, ok := err.(awserr.Error); ok {
      switch aerr.Code() {
      default:
        fmt.Println(aerr.Error())
      }
    } else {
      // Print the error, cast err to awserr.Error to get the Code and
      // Message from an error.
      fmt.Println(err.Error())
    }
    return
  }
  fmt.Println(result)
}

//To GET Post body from Table "Post" in DynamoDB by "post_id". Return Post body
func GetPostBody(post Post, sess *session.Session) Post {
  svc := dynamodb.New(sess)

  //Request to DynamoDB to GET post by "post_id"
  input := &dynamodb.GetItemInput{
    Key: map[string]*dynamodb.AttributeValue{
      "post_id": {
        S: &post.PostID,
      },
    },
    TableName: aws.String("Post"),
  }

  //Get result
  result, err := svc.GetItem(input)

  //ERROR block (from AWS SDK GO DynamoDB documentation)
  if err != nil {
    fmt.Println(err)
  }

  //Unmarshal result to Post structure
  err = dynamodbattribute.UnmarshalMap(result.Item, &post)

  return post
}

//To generate a random UUID according to RFC 4122
func newUUID() (string, error) {
  uuid := make([]byte, 16)
  n, err := io.ReadFull(rand.Reader, uuid)
  if n != len(uuid) || err != nil {
    return "", err
  }
  // variant bits; see section 4.1.1
  uuid[8] = uuid[8]&^0xc0 | 0x80
  // version 4 (pseudo-random); see section 4.1.3
  uuid[6] = uuid[6]&^0xf0 | 0x40
  return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}


