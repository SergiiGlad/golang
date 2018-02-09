package controllers

import (
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
  "github.com/aws/aws-sdk-go/aws/awserr"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "fmt"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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

func GetPost (svc dynamodbiface.DynamoDBAPI) (Post, error){
  var post Post
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

}
