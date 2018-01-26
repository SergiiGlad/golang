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

//Post structure
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

  //Decode request JSON
  _ = json.NewDecoder(r.Body).Decode(&newPost)

  //Create new Session for DynamoDB
  sess, err := session.NewSession(&aws.Config{
    Region: aws.String("eu-west-2"),
  })
  svc := dynamodb.New(sess)

  //Request to DynamoDB to CREATE new post with KEY_ATTRIBUTE "post_id" (TimeStamp)
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

//To GET posts from DynamoDB Table "Post" by UserID
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
  proj := expression.NamesList(expression.Name("post_title"), expression.Name("post_text"), expression.Name("post_id"), expression.Name("user_id"))

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

