package messages

import (
	"testing"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/gusaul/go-dynamock"
)

func TestPutMessageToDynamo(t *testing.T) {
	var mock *dynamock.DynaMock
	Dyna.Db, mock = dynamock.New()
	///////////////
	expectKey := map[string]*dynamodb.AttributeValue{
		"id": {
			N: aws.String("23"),
		},
	}

	expectedResult := aws.String("Vasya")
	result := dynamodb.GetItemOutput{
		Item: map[string]*dynamodb.AttributeValue{
			"name": {
				S: expectedResult,
			},
		},
	}

	//lets start dynamock in action
	mock.ExpectGetItem().ToTable("massages").WithKeys(expectKey).WillReturns(result)

}
