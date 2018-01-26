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
  "github.com/gorilla/mux"
  "github.com/aws/aws-sdk-go/service/dynamodb/expression"
  "os"
)

type Post struct{
  Title string `json:"post_title"`
  Text string	`json:"post_text"`
  PostID string `json:"post_id"`
  UserID string `json:"user_id"`
  Like string `json:"post_like"`
}

//To CREATE new post in DynamoDB Table "Post"
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
      "user_id": {
        S: &newPost.UserID,
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

  _ = json.NewEncoder(w).Encode(&newPost)

  fmt.Println(result)
}

//To DELETE existing post by "post_id" from DynamoDB Table "Post"
func DeletePost(w http.ResponseWriter, r *http.Request) {

  var post Post
  vars := mux.Vars(r)
  post.PostID = vars["post_id"]

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

  _ = json.NewEncoder(w).Encode(&post.PostID)

  fmt.Println(result)
}

//To GET post by "post_id" from DynamoDB Table "Post"
func GetPost(w http.ResponseWriter, r *http.Request) {
  var post Post

  vars := mux.Vars(r)
  post.PostID = vars["post_id"]

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

  fmt.Println("ID: ", post.PostID,"User ID: ", post.UserID, "Title: ", post.Title, "Text: ", post.Text)
  _ = json.NewEncoder(w).Encode(&post)
}

//To GET posts from DynamoDB Table "Post" by UserID
func GetPostByUserID(w http.ResponseWriter, r *http.Request){
  var post Post

  vars := mux.Vars(r)
  post.UserID = vars["user_id"]

  sess, err := session.NewSession(&aws.Config{
    Region:      aws.String("eu-west-2"),
  })
  svc := dynamodb.New(sess)

  filt := expression.Name("user_id").Equal(expression.Value(post.UserID))

  proj := expression.NamesList(expression.Name("post_title"), expression.Name("post_text"), expression.Name("post_id"), expression.Name("user_id"))

  expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

  params := &dynamodb.ScanInput{
    ExpressionAttributeNames:  expr.Names(),
    ExpressionAttributeValues: expr.Values(),
    FilterExpression:          expr.Filter(),
    ProjectionExpression:      expr.Projection(),
    TableName:                 aws.String("Post"),
  }


  // Make the DynamoDB Query API call
  result, err := svc.Scan(params)

  for _, i := range result.Items {
    post := Post{}

    err = dynamodbattribute.UnmarshalMap(i, &post)

    if err != nil {
      fmt.Println("Got error unmarshalling:")
      fmt.Println(err.Error())
      os.Exit(1)
    }

    _ = json.NewEncoder(w).Encode(&post)

    fmt.Println("Post ID: ", post.PostID)
    fmt.Println("Post Title:", post.Title)
    fmt.Println("Post Text:", post.Text)
    fmt.Println()
  }
}

//To UPDATE post in DynamoDB Table "Post"
func UpdatePost(w http.ResponseWriter, r *http.Request){
  var post Post


  _ = json.NewDecoder(r.Body).Decode(&post)

  vars := mux.Vars(r)
  post.PostID = vars["post_id"]
  fmt.Println(post.PostID)

  sess, err := session.NewSession(&aws.Config{
    Region: aws.String("eu-west-2"),
  })

  svc := dynamodb.New(sess)
  input := &dynamodb.UpdateItemInput{
    ExpressionAttributeNames: map[string]*string{
      "#PTitle": aws.String("post_title"),
      "#PText": aws.String("post_text"),
    },
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

  result, err := svc.UpdateItem(input)


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

  _ = json.NewEncoder(w).Encode(&post)

}

//To describe table "Post" in DynamoDB
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

  _ = json.NewEncoder(w).Encode(result)
  fmt.Println(result)
}

