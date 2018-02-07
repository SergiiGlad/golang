package messages

import (
	"encoding/json"
	"fmt"
	"go-team-room/conf"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

//GetMessageFromDynamoByUserId - return an slise of latest
//messages but not more then MAX_MESSAGES_AT_ONCE
func GetMessageFromDynamoByUserID(humUserId int, maxMessages ...int) []HumMessage {
	var resultMessages []HumMessage
	maxMes := conf.MaxMessages
	if maxMessages != nil && maxMessages[0] < maxMes {
		maxMes = maxMessages[0]
	}
	if humUserId < 1 {
		return resultMessages //emty array
	}
	// Create the Expression to fill the input struct with.
	filt := expression.Name("message_user.id_sql").Equal(expression.Value(humUserId))
	expr, err := expression.NewBuilder().WithFilter(filt).Build()

	if err != nil {
		fmt.Println("Got error building Dynamo expression:")
		fmt.Println(err.Error()) //DEBUG output
		//os.Exit(1)
		return resultMessages //emty array
	}
	// Build the query input parameters
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		//  ProjectionExpression:      expr.Projection(),
		TableName: aws.String("messages"),
	}
	result, err := Dyna.Db.Scan(params)
	if err != nil {
		fmt.Println("Query API call failed:")
		fmt.Println((err.Error()))
		//fmt.Println(params)
		//os.Exit(1)
		return resultMessages //empty
	}
	var humMessageHolder HumMessage
	for _, i := range result.Items {

		//err := json.Unmarshal(i, humMessageHolder)
		//err := json.NewDecoder(i).Decode(&humMessageHolder)
		err = dynamodbattribute.UnmarshalMap(i, &humMessageHolder)
		if err != nil {
			fmt.Println(err) //panic(err)
		} else {
			resultMessages = append(resultMessages, humMessageHolder)
		}

		//fmt.Println("message_id: ", item.MessageId)
		//fmt.Println("Message Data Text:", item.MessageData.Text)
	}
	return resultMessages
}

//HandlerOfGetMessages - handler for
//"/messages/" endpoint

func HandlerOfGetMessages(writeRespon http.ResponseWriter, r *http.Request) {
	currentUserID := GetActionUserID(r)
	//fmt.Println(currentUserID)
	if currentUserID < 1 {
		fmt.Fprint(writeRespon, "403 Can't find your ID")
		return
	}
	var numberOfMessages int
	keys, ok := r.URL.Query()["numberOfMessages"]
	if !ok || len(keys) < 1 {
		numberOfMessages = conf.MaxMessages
	} else {
		numberOfMessages, _ = strconv.Atoi(keys[0])
		if numberOfMessages > conf.MaxMessages || numberOfMessages < 1 {
			numberOfMessages = conf.MaxMessages
		}
	}
	usersMessagesObj := GetMessageFromDynamoByUserID(currentUserID, numberOfMessages)
	var resultString string

	if len(usersMessagesObj) > 0 {
		for i := 0; i < len(usersMessagesObj); i++ {
			tempMarshaledMessage, err := json.Marshal(usersMessagesObj[i])
			if err != nil {
				//panic(err)
				//it is impossible!!!
				//we are iterating via Hum obj
			} else {
				resultString = resultString + string(tempMarshaledMessage)
			}
		}
	}
	fmt.Fprint(writeRespon, resultString)
	return
}
