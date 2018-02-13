package controllers

import (
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
  "github.com/aws/aws-sdk-go/aws/awserr"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "fmt"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
  "github.com/aws/aws-sdk-go/service/dynamodb/expression"
  "os"
  "github.com/aws/aws-sdk-go/aws/session"
  "mime/multipart"
  "strings"
  "github.com/aws/aws-sdk-go/service/s3/s3manager"
  "go-team-room/conf"
  "io"
  "crypto/rand"
  "github.com/pkg/errors"
  "github.com/aws/aws-sdk-go/service/s3/s3iface"
  "github.com/aws/aws-sdk-go/service/s3"
)

//Post structure
type Post struct{
  Title string `json:"post_title"`
  Text string	`json:"post_text"`
  PostID string `json:"post_id"`
  UserID string `json:"user_id"`
  Like string `json:"post_like"`
  FileLink string `json:"file_link"`
  LastUpdate string `json:"post_last_update"`
}


func CreateNewPost(svc dynamodbiface.DynamoDBAPI, post Post) (Post, error)  {


  //Request to DynamoDB to CREATE new post with KEY_ATTRIBUTE "post_id" (TimeStamp)
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
        N: &post.Like,
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
    return post, err
  }
  fmt.Println(result)

  return post, nil
}

func GetPost(svc dynamodbiface.DynamoDBAPI, post_id string) (Post, error){
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
  }

  //Check if POST exists in table, if not: RESPONSE 204 - "No Content"
  if len(result.Item) == 0 {
    return post, errors.New("No content")
  }

  //Unmarshal result to Post structure
  err = dynamodbattribute.UnmarshalMap(result.Item, &post)

  //Print result in console
  fmt.Println(result)
  return post, nil
}

func GetPostByUserID(svc dynamodbiface.DynamoDBAPI, user_id string) ([]Post, error){
  var post Post
  var outputPost []Post

  post.UserID = user_id

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
    return outputPost, err
  }

  //Loop for encoding multiples post which was made by "user_id"
  for _, i := range result.Items {
    post := Post{}
    //Unmarshal result to Post structure
    err = dynamodbattribute.UnmarshalMap(i, &post)
    outputPost = append(outputPost, post)

    //Check error for unmarshalling
    if err != nil {
      fmt.Println("Got error unmarshalling:")
      fmt.Println(err.Error())
      os.Exit(1)
    }
  }
  //Print result in console
  fmt.Print(result)
  return outputPost, nil
}

func UpdatePost(svc dynamodbiface.DynamoDBAPI, post Post) (Post, error){

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
    return post, err
  }

  //Print result in console
  fmt.Println(result)
  return post, nil
}

func DeletePost(svcd dynamodbiface.DynamoDBAPI, svcs s3iface.S3API, post_id string) string{

  post, err := GetPost(svcd, post_id)

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
    return "Error"
  }

  fmt.Println(result)

  return post_id
}

//To UPLOAD file to S3
func UploadFileToS3(sess *session.Session, f multipart.File, handl * multipart.FileHeader) string{
  // Create an uploader with the session and default options
  uploader := s3manager.NewUploader(sess)

  fileType := handl.Filename[strings.LastIndexAny(handl.Filename, "."):]
  fileType = strings.ToLower(fileType)

  //Generate UUID for File
  uuid, err := newUUID()
  if err != nil{
    fmt.Printf("error: %v\n", err)
  }

  // Upload the file to S3.
  result, err := uploader.Upload(&s3manager.UploadInput{
    Bucket: aws.String(conf.AwsBucketName),
    Key:    aws.String(uuid + fileType),
    Body:   f,
  })
  if err != nil{
    fmt.Errorf("failed to upload file, %v", err)
    return "failed to upload file"
  }
  fmt.Printf("file uploaded to, %s\n", aws.StringValue(&result.Location))
  link := "/uploads/" + uuid + fileType
  return link
}

//To DELETE file from S3 when DELETE Post
func DeleteFileFromS3(file_link string, svc s3iface.S3API) {
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
