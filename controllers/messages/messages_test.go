package messages

import (
	"encoding/json"
	"fmt"
	"go-team-room/conf"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

func TestGetMessagesFromDynamoByID  (t *testing.T){
	var mock *dynamock.DynaMock
	Dyna.db, mock = dynamock.New()
	///////////////
	expectKey := map[string]*dynamodb.AttributeValue{
        "id" : {
            N: aws.String("1")
        },
    }

    expectedResult := aws.String("qwer")
    result := dynamodb.GetItemOutput{
        Item: map[string]*dynamodb.AttributeValue{
            "name": {
                S: expectedResult,
            },
        },
    }

    // start dynamock in action
    mock.ExpectGetItem().ToTable("employee").WithKeys(expectKey).WillReturns(result)

}