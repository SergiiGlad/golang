package server

import (
  "net/http"
  "encoding/json"
  "fmt"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "time"
  "github.com/aws/aws-sdk-go/aws/awserr"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Post struct{
  Title string `json:"post_title"`
  Text string	`json:"post_text"`
  PostID string `json:"post_id"`
}

//To create new post in DynamoDB Table "Post"
func CreateNewPost(w http.ResponseWriter, r *http.Request) {
  var newPost Post

  _ = json.NewDecoder(r.Body).Decode(&newPost)

  fmt.Println(newPost)

  sess, err := session.NewSession(&aws.Config{
    Region: aws.String("eu-west-2"),
  })

  svc := dynamodb.New(sess)
  input := &dynamodb.PutItemInput{
    Item: map[string]*dynamodb.AttributeValue{
      "post_id": {
        S: aws.String(time.Now().String()),
      },
      "post_title": {
        S: &newPost.Title,
      },
      "post_text": {
        S: &newPost.Text,
      },
    },
    ReturnConsumedCapacity: aws.String("TOTAL"),
    TableName:              aws.String("Post"),
  }

  result, err := svc.PutItem(input)
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

  fmt.Println(result)
}

//To delete existing post by "post_id" from DynamoDB Table "Post"
func DeletePost(w http.ResponseWriter, r *http.Request) {

  var post Post

  _ = json.NewDecoder(r.Body).Decode(&post)
  fmt.Println(post.PostID)

  sess, err := session.NewSession(&aws.Config{
    Region:      aws.String("eu-west-2"),
  })

  svc := dynamodb.New(sess)

  input := &dynamodb.DeleteItemInput{
    Key: map[string]*dynamodb.AttributeValue{
      "post_id": {
        S: &post.PostID,
      },
    },
    TableName: aws.String("Post"),
  }

  result, err := svc.DeleteItem(input)
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

  fmt.Println(result)
}

//To get post by "post_id" from DynamoDB Table "Post"
func GetPost(w http.ResponseWriter, r *http.Request) {
  var post Post
  _ = json.NewDecoder(r.Body).Decode(&post)

  sess, err := session.NewSession(&aws.Config{
    Region:      aws.String("eu-west-2"),
  })
  svc := dynamodb.New(sess)
  input := &dynamodb.GetItemInput{
    Key: map[string]*dynamodb.AttributeValue{
      "post_id": {
        S: &post.PostID,
      },
    },
    TableName: aws.String("Post"),
  }

  result, err := svc.GetItem(input)

  err = dynamodbattribute.UnmarshalMap(result.Item, &post)

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

  fmt.Println(result)

  fmt.Println("ID: ", post.PostID, "Title: ", post.Title, "Text: ", post.Text)
  _ = json.NewEncoder(w).Encode(&post)
}

//To decribe table "Post" in DynamoDB
func DescribeTablePost(w http.ResponseWriter, r *http.Request) {

  sess, err := session.NewSession(&aws.Config{
    Region:      aws.String("eu-west-2"),
  })
  svc := dynamodb.New(sess)
  input := &dynamodb.DescribeTableInput{
    TableName: aws.String("Post"),
  }

  result, err := svc.DescribeTable(input)
  if err != nil {
    if aerr, ok := err.(awserr.Error); ok {
      switch aerr.Code() {
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

  fmt.Println(result)
}
