package Amazon

import (
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "go-team-room/conf"
  "fmt"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)
var SVC *dynamodb.DynamoDB
var SESS *session.Session

type MyDynamo struct {
  Db dynamodbiface.DynamoDBAPI
}

var Dynamo MyDynamo

func init() {
  //Create Session for AWS
    SESS, err := session.NewSession(&aws.Config{
    Region: aws.String(conf.DynamoRegion),
  })
    if err != nil {
      fmt.Println(err)
    }
  SVC = dynamodb.New(SESS)
  Dynamo.Db = dynamodbiface.DynamoDBAPI(SVC)
}
