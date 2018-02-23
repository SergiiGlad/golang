package controllers

import (
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
  "github.com/aws/aws-sdk-go/aws/awserr"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "fmt"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
  "github.com/aws/aws-sdk-go/service/dynamodb/expression"
  "mime/multipart"
  "strings"
  "github.com/aws/aws-sdk-go/service/s3/s3manager"
  "go-team-room/conf"
  "io"
  "crypto/rand"
  "github.com/pkg/errors"
  "github.com/aws/aws-sdk-go/service/s3/s3iface"
  "github.com/aws/aws-sdk-go/service/s3"
  "go-team-room/humstat"
)

//Post structure
type Post struct {
  Title  string `json:"post_title"`
  Text   string `json:"post_text"`
  PostID string `json:"post_id"`
  UserID string  `json:"user_id"`
  //Like string `json:"post_like"`
  Like       [] *string `json:"post_like"`
  FileLink   string   `json:"file_link"`
  LastUpdate string   `json:"post_last_update"`
}

//To CREATE new Post in DynamoDB
func CreateNewPost(svc dynamodbiface.DynamoDBAPI, post Post) (Post, error) {

  //Request to DynamoDB to CREATE new post with KEY_ATTRIBUTE "post_id"
  input := &dynamodb.PutItemInput{
    Item: map[string]*dynamodb.AttributeValue{
      "post_id": {
        S: &post.PostID,
      },
      "post_title": {
        S: &post.Title,
      },
      "post_text": {
        S: &post.Text,
      },
      "user_id": {
        S: &post.UserID,
      },
      "post_like": {
        SS: post.Like,
      },
      "file_link": {
        S: &post.FileLink,
      },
      "post_last_update": {
        S: &post.LastUpdate,
      },
    },
    ReturnConsumedCapacity: aws.String("TOTAL"),
    TableName:              aws.String("Post"), //Name of Table in DynamoDB
  }

  //Get result
  result, err := svc.PutItem(input)

  //ERROR block
  if err != nil {
    if aerr, ok := err.(awserr.Error); ok {
      switch aerr.Code() {
      case dynamodb.ErrCodeConditionalCheckFailedException:
        log.Error(dynamodb.ErrCodeConditionalCheckFailedException, aerr.Error())
      case dynamodb.ErrCodeProvisionedThroughputExceededException:
        log.Error(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
      case dynamodb.ErrCodeResourceNotFoundException:
        log.Error(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
      case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
        log.Error(dynamodb.ErrCodeItemCollectionSizeLimitExceededException, aerr.Error())
      case dynamodb.ErrCodeInternalServerError:
        log.Error(dynamodb.ErrCodeInternalServerError, aerr.Error())
      default:
        log.Error(aerr.Error())
      }
    } else {
      // Print the error, cast err to awserr.Error to get the Code and
      // Message from an error.
      log.Error(err.Error())
    }
    return post, errors.New("Error")
  }
  _ = result

  //logging
  log.Info("Post created successfully")

  //statistic
  humstat.SendStat <- map[string]int {
    "Post created": 1,
  }

  return post, nil
}

//To GET Post by POST_ID from DynamoDB
func GetPost(svc dynamodbiface.DynamoDBAPI, post_id string) (Post, error) {
  var post Post

  post.PostID = post_id

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

  //ERROR block
  if err != nil {
    if aerr, ok := err.(awserr.Error); ok {
      switch aerr.Code() {
      case dynamodb.ErrCodeProvisionedThroughputExceededException:
        log.Error(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
      case dynamodb.ErrCodeResourceNotFoundException:
        log.Error(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
      case dynamodb.ErrCodeInternalServerError:
        log.Error(dynamodb.ErrCodeInternalServerError, aerr.Error())
      default:
        log.Error(aerr.Error())
      }
    } else {
      // Print the error, cast err to awserr.Error to get the Code and
      // Message from an error.
      log.Error(err.Error())
    }
  }

  //Check if POST exists in table, if not: RESPONSE 204 - "No Content"
  if len(result.Item) == 0 {
    log.Info("No content")
    return post, errors.New("No content")
  }

  //Unmarshal result to Post structure
  err = dynamodbattribute.UnmarshalMap(result.Item, &post)

  if err != nil {
    log.Error("Error to unmarshal")
    return post, errors.New("Error to unmarshal")
  }

  //logging
  log.Info("Get post by post id successfully")

  //statistic
  humstat.SendStat <- map[string]int{
    "Get post by post id": 1,
  }

  return post, nil
}

//To GET Posts by USER_ID from DynamoDB
func GetPostByUserID(svc dynamodbiface.DynamoDBAPI, user_id string) ([]Post, error) {

  var outputPost []Post

  //Filter expression: Seeks all items in table with equal "user_id"
  filt := expression.Name("user_id").Equal(expression.Value(user_id))

  //Make projection: displays all expression.Name with equal "user_id"
  proj := expression.NamesList(expression.Name("post_title"), expression.Name("post_text"), expression.Name("post_id"), expression.Name("user_id"), expression.Name("post_like"), expression.Name("file_link"), expression.Name("post_last_update"))

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
    return outputPost, errors.New("No content")
  }

  //Loop for encoding multiples post which was made by "user_id"
  for _, i := range result.Items {
    post := Post{}

    //Unmarshal result to Post structure
    err = dynamodbattribute.UnmarshalMap(i, &post)

    //Check error for unmarshalling
    if err != nil {
      log.Error("Got error unmarshalling:")
      return outputPost, errors.New("Error unmarshalling")
    }

    outputPost = append(outputPost, post)
  }

  //logging
  log.Info("Get post by user id successfully")

  //statistic
  humstat.SendStat <- map[string]int {
    "Get post by user id": 1,
  }

  return outputPost, nil
}

//To UPDATE Post TITLE, TEXT, POST_LAST_UPDATE in DynamoDB
func UpdatePost(svc dynamodbiface.DynamoDBAPI, post Post) (Post, error) {

  //Request to DynamoDB to UPDATE Item in table
  input := &dynamodb.UpdateItemInput{
    ExpressionAttributeNames: map[string]*string{
      "#PTitle": aws.String("post_title"),
      "#PText":  aws.String("post_text"),
      "#PDate":  aws.String("post_last_update"),
    },
    //Attributes to UPDATE
    ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
      ":t": {
        S: &post.Title,
      },
      ":e": {
        S: &post.Text,
      },
      ":d": {
        S: &post.LastUpdate,
      },
    },
    Key: map[string]*dynamodb.AttributeValue{
      "post_id": {
        S: &post.PostID,
      },
    },
    ReturnValues:     aws.String("UPDATED_NEW"),
    TableName:        aws.String("Post"),
    UpdateExpression: aws.String("set #PTitle = :t, #PText = :e, #PDate = :d"),
  }

  //Get result
  result, err := svc.UpdateItem(input)

  //ERROR block (from AWS SDK GO DynamoDB documentation)
  if err != nil {
    if aerr, ok := err.(awserr.Error); ok {
      switch aerr.Code() {
      case dynamodb.ErrCodeConditionalCheckFailedException:
        log.Error(dynamodb.ErrCodeConditionalCheckFailedException, aerr.Error())
      case dynamodb.ErrCodeProvisionedThroughputExceededException:
        log.Error(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
      case dynamodb.ErrCodeResourceNotFoundException:
        log.Error(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
      case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
        log.Error(dynamodb.ErrCodeItemCollectionSizeLimitExceededException, aerr.Error())
      case dynamodb.ErrCodeInternalServerError:
        log.Error(dynamodb.ErrCodeInternalServerError, aerr.Error())
      default:
        log.Error(aerr.Error())
      }
    } else {
      // Print the error, cast err to awserr.Error to get the Code and
      // Message from an error.
      log.Error(err.Error())
    }
    return post, errors.New("Error")
  }

  _ = result

  //logging
  log.Info("Update post successfully")

  //statistic
  humstat.SendStat <- map[string]int {
    "Update post": 1,
  }

  return post, nil
}

//To DELETE Post in DynamoDB
func DeletePost(svcd dynamodbiface.DynamoDBAPI, svcs s3iface.S3API, post Post) string{

  //Check if Post has file on S3
  //if true -> Delete file from S3
  if post.FileLink != "NULL" {
    DeleteFileFromS3(post.FileLink, svcs)
  }

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
  result, err := svcd.DeleteItem(input)

  //ERROR block (from AWS SDK GO DynamoDB documentation)
  if err != nil {
    if aerr, ok := err.(awserr.Error); ok {
      switch aerr.Code() {
      case dynamodb.ErrCodeConditionalCheckFailedException:
        log.Error(dynamodb.ErrCodeConditionalCheckFailedException, aerr.Error())
      case dynamodb.ErrCodeProvisionedThroughputExceededException:
        log.Error(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
      case dynamodb.ErrCodeResourceNotFoundException:
        log.Error(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
      case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
        log.Error(dynamodb.ErrCodeItemCollectionSizeLimitExceededException, aerr.Error())
      case dynamodb.ErrCodeInternalServerError:
        log.Error(dynamodb.ErrCodeInternalServerError, aerr.Error())
      default:
        log.Error(aerr.Error())
      }
    } else {
      // Print the error, cast err to awserr.Error to get the Code and
      // Message from an error.
      log.Error(err.Error())
    }
  }

  _ = result

  //logging
  log.Info("Delete post successfully")

  //statistic
  humstat.SendStat <- map[string]int {
    "Delete post": 1,
  }

  return "Deleted"
}

//To UPLOAD file to S3
func UploadFileToS3(svc s3iface.S3API, f multipart.File, handl *multipart.FileHeader) string {
  // Create an uploader with the session and default options
  uploader := s3manager.NewUploaderWithClient(svc)

  //Get type of file
  fileType := handl.Filename[strings.LastIndexAny(handl.Filename, "."):]
  fileType = strings.ToLower(fileType)

  //Generate UUID for File
  uuid, err := NewUUID()
  if err != nil {
    fmt.Printf("error: %v\n", err)
  }

  // Upload the file to S3.
  result, err := uploader.Upload(&s3manager.UploadInput{
    Bucket: aws.String(conf.AwsBucketName),
    Key:    aws.String(uuid + fileType),
    Body:   f,
  })

  if err != nil {
    fmt.Errorf("failed to upload file, %v", err)
    return "failed to upload file"
  }

  _ = result

  //logging
  log.Info("File uploaded successfully")

  //statistic
  humstat.SendStat <- map[string]int {
    "File upload": 1,
  }

  //Get link to store it in DynamoDB
  link := "/uploads/" + uuid + fileType

  return link
}

//To DELETE file from S3 when DELETE Post
func DeleteFileFromS3(file_link string, svc s3iface.S3API) {
  //Get file name from DynamoDB, trimming "/uploads/"
  filename := strings.TrimPrefix(file_link, "/uploads/")

  input := &s3.DeleteObjectsInput{
    Bucket: aws.String(conf.AwsBucketName),
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
        log.Error(aerr.Error())
      }
    } else {
      // Print the error, cast err to awserr.Error to get the Code and
      // Message from an error.
      log.Error(err.Error())
    }
    return
  }
  _ = result

  //logging
  log.Info("File deleted")

  //statistic
  humstat.SendStat <- map[string]int {
    "File deleted": 1,
  }
}

//To DOWNLOAD file from S3
func DownloadFileFromS3(svc s3iface.S3API, fileName string) (*aws.WriteAtBuffer, error) {
  // Create a downloader with the session and default options
  downloader := s3manager.NewDownloaderWithClient(svc)

  var b []byte
  buff := aws.NewWriteAtBuffer(b)

  // Write the contents of S3 Object to the file
  n, err := downloader.Download(buff, &s3.GetObjectInput{
    Bucket: aws.String(conf.AwsBucketName),
    Key:    aws.String(fileName),
  })

  if err != nil {
    fmt.Errorf("failed to download file, %v", err)
    return buff, errors.New("Error")
  }

  //logging
  log.Info("File downloaded, %d bytes\n", n)

  //statistic
  humstat.SendStat <- map[string]int {
    "File downloaded": 1,
  }

  return buff, nil
}

//To generate a random UUID according to RFC 4122
func NewUUID() (string, error) {
  uuid := make([]byte, 16)
  n, err := io.ReadFull(rand.Reader, uuid)
  if n != len(uuid) || err != nil {
    return "", err
  }
  // variant bits; see section 4.1.1
  uuid[8] = uuid[8]&^0xc0 | 0x80
  // version 4 (pseudo-random); see section 4.1.3
  uuid[6] = uuid[6]&^0xf0 | 0x40
  log.Info("New UUID has generated")
  return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
