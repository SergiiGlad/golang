package messages

import (
	"fmt"
	"go-team-room/conf"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type MyDynamo struct {
	Db dynamodbiface.DynamoDBAPI
}

var Dyna *MyDynamo

func init() {
	Dyna = new(MyDynamo)
	awsSession, err := session.NewSession(&aws.Config{
		Region:      aws.String(conf.DynamoRegion),
		Credentials: credentials.NewStaticCredentials(conf.AwsAccessKeyId, conf.AwsSecretKey, ""),
	})
	if err != nil {
		fmt.Println("Got error creating Dynamo session (can't init messages)")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var svc *dynamodb.DynamoDB = dynamodb.New(awsSession)
	Dyna.Db = dynamodbiface.DynamoDBAPI(svc)

}
