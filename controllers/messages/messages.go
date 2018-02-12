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

//MyDynamo struct for mock`ing DynampDB in tests
type MyDynamo struct {
	Db dynamodbiface.DynamoDBAPI
}

//Dyna - global variable to use as DynamoDb interfase
var Dyna *MyDynamo

// Get instance of logger (Formatter, Hook，Level, Output ).
// If you want to use only your log message  It will need use own call logs example
var logRus = conf.GetLog()

func init() {
	//logRus = logRus.WithField(logRus.Fields{"packet": "messages"})

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
	logRus.Debug("Messages packet initialized Ok")
}
