package Amazon

import (
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "go-team-room/conf"
  "fmt"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
  "github.com/aws/aws-sdk-go/service/s3/s3iface"
  "github.com/aws/aws-sdk-go/service/s3"
)

var SVCD *dynamodb.DynamoDB
var SVCS *s3.S3


type MyDynamo struct {
  Db dynamodbiface.DynamoDBAPI
}

type MyS3 struct{
  S3API s3iface.S3API
}

var Dynamo MyDynamo
var S3 MyS3

func init() {
  //Create Session for AWS
    SESS, err := session.NewSession(&aws.Config{
    Region: aws.String(conf.DynamoRegion),
  })
    if err != nil {
      fmt.Println(err)
    }
  SVCD = dynamodb.New(SESS)
  Dynamo.Db = dynamodbiface.DynamoDBAPI(SVCD)
  SVCS = s3.New(SESS)
  S3.S3API = s3iface.S3API(SVCS)
}
