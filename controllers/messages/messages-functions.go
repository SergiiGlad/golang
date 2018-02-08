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

//GetMessageFromDynamoByUserID - return an slise of latest
//messages but not more then MAX_MESSAGES_AT_ONCE
func GetMessageFromDynamoByUserID(humUserID int, maxMessages ...int) []HumMessage {
	var resultMessages []HumMessage
	maxMes := conf.MaxMessages
	if maxMessages != nil && maxMessages[0] < maxMes {
		maxMes = maxMessages[0]
	}
	if humUserID < 1 {
		return resultMessages //emty array
	}
	// Create the Expression to fill the input struct with.
	filt := expression.Name("message_user.id_sql").Equal(expression.Value(humUserID))
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
//"/messages" endpoint
func HandlerOfGetMessages(writeRespon http.ResponseWriter, r *http.Request) {
	///Debug
	GetChatRoomListByUserID(23, 23)
	return
	///DEBUG

	currentUserID := GetActionUserID(r)
	//fmt.Println(currentUserID)
	if currentUserID < 1 {
		writeRespon.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(writeRespon, "401 Can't find your ID")
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
			//fmt.Println(tempMarshaledMessage)
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

//GetActionUserID - return ID of user who make this action
func GetActionUserID(r *http.Request) int {
	// TEMPORARY stub
	//Later it will get user from session/token/cookies
	iid := -1
	if r.Method == "GET" {

		keys, ok := r.URL.Query()["id"]
		if !ok || len(keys) < 1 {
			iid = -1
		} else {
			iiid, err := strconv.Atoi(keys[0])
			if err != nil {
				iid = -1
			} else {
				iid = iiid
			}
		}
	} else if r.Method == "POST" {
		var data HumMessage
		fmt.Println(r.Body)
		err := json.NewDecoder(r.Body).Decode(&data)
		fmt.Println(r.Body)
		if err != nil {
			//panic(err)
			iid = -1
		} else {
			iid = data.MessageUser.IdSql
		}
	} else {
		iid = -1 //UNsupported method/bad ID
	}
	return iid

}

//ValidateDataFromUser very important func
//Do NOT trust any data from User!!!
//VALIDATE EVERETHING
func ValidateDataFromUser(m *HumMessage) {
	//work this out LATER
	//for now it is stub only
	//
	//do nothing
	//go home
}

//ReadReqBodyPOST - reads all POST request body to HumMessage
func ReadReqBodyPOST(req *http.Request, humMess *HumMessage) {
	// body, err := ioutil.ReadAll(req.Body)
	//  if err != nil {
	//  	return //r1
	//  }
	// fmt.Println(string(body)) //DEBUG output
	var data HumMessage

	err := json.NewDecoder(req.Body).Decode(&data)
	fmt.Println(data)
	//err = json.Unmarshal(body, humMess)
	if err != nil {
		//panic(err)
		fmt.Println(err)
	}
}

//PutMessageToDynamo get prepeared HummMessage obj and write it to Dynamo
//result of operation writed into writeRespon
func PutMessageToDynamo(writeRespon http.ResponseWriter, m *HumMessage) {

	av, err := dynamodbattribute.MarshalMap(m)

	if err != nil {
		fmt.Println("Got error marshalling map:")
		fmt.Println(err.Error())
		//os.Exit(1)
	}

	// Create item in table messages
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("messages"),
	}

	_, err = Dyna.Db.PutItem(input)

	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		//os.Exit(1)
		writeRespon.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(writeRespon, "Some errors")
	} else {
		writeRespon.WriteHeader(http.StatusOK)

		fmt.Fprint(writeRespon, "Post done")
	}
	//fmt.Println("Successfully added 'The Big someNewMessage' to  table")
}

//HandlerOfPOSTMessages This func should process any messages end point
func HandlerOfPOSTMessages(w http.ResponseWriter, r *http.Request) {

	//DEBUG
	//currentUserID := GetActionUserID(r)
	var inputMessage HumMessage
	err := json.NewDecoder(r.Body).Decode(&inputMessage)
	if err != nil {
		//panic(err)
		return
	}

	//CORS!!! "C:\Program Files (x86)\Google\Chrome\Application\chrome.exe" --disable-web-security --user-data-dir="D:/Chrome"

	//Assume it is an POST
	//Shuold  expect a user message in form of
	/*
		 curl -X POST --header 'Content-Type: application/json' --header 'Accept: application/json' -d '{"message_chat_room_id": "997",
		   "message_data": {
			 "binary_parts": [{
		 "bin_data": null,
		  "bin_name": null }],
			 "text": "A lot of text and stupid smiles :)))))",
		  "type": "TypeOfHumMessage-UNDEFINED FOR NOW"},
		   "message_id": "20180110155343152",
		   "message_parent_id": "20180110155533289",
		   "message_social_status": {
			 "Dislike": 11,
			 "Like": 22,
			 "Views": 33 },
		   "message_timestamp": "20180110155533111",
		   "message_user": {
			 "id_sql": 23,
			 "name_sql": "Vasya" }
		 }' 'http://localhost:8080/messages'
	*/
	//// https://gist.github.com/alyssaq/75d6678d00572d103106

	//ReadReqBodyPOST(r, &inputMessage) //Do not use it NOW.

	ValidateDataFromUser(&inputMessage) //Its FAKE

	//assume data validated
	//and it is safe to put it into a Dynamo

	PutMessageToDynamo(w, &inputMessage)
}

func GetChatRoomListByUserID(humUserID int, maxChatRooms ...int) []HumChatRoom {
	var resultChatRooms []HumChatRoom
	maxChats := conf.MaxChatRooms
	if maxChatRooms != nil && maxChatRooms[0] < maxChats {
		maxChats = maxChatRooms[0]
	}
	if humUserID < 1 {
		return resultChatRooms //emty array
	}
	// Create the Expression to fill the input struct with.
	filt := expression.Name("chat_users_list[0].id_sql").Equal(expression.Value(humUserID))
	expr, err := expression.NewBuilder().WithFilter(filt).Build()

	if err != nil {
		fmt.Println("Got error building Dynamo expression:")
		fmt.Println(err.Error()) //DEBUG output
		//os.Exit(1)
		return resultChatRooms //emty array
	}
	// Build the query input parameters
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		//  ProjectionExpression:      expr.Projection(),
		TableName: aws.String("chat_rooms"),
	}
	result, err := Dyna.Db.Scan(params)
	if err != nil {
		fmt.Println("Query API call failed:")
		fmt.Println((err.Error()))
		//fmt.Println(params)
		//os.Exit(1)
		return resultChatRooms //empty
	}
	var humChatHolder HumChatRoom
	for _, i := range result.Items {
		err = dynamodbattribute.UnmarshalMap(i, &humChatHolder)
		if err != nil {
			fmt.Println(err) //panic(err)
		} else {
			resultChatRooms = append(resultChatRooms, humChatHolder)
		}
		fmt.Println("item.ChatID: ", i)
		//fmt.Println("item.ChatUsersList", i.ChatUsersList)
	}
	return resultChatRooms
}
