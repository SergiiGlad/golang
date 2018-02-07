package messages

import (
	"encoding/json"
	"fmt"
	"go-team-room/conf"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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
		//decoder := json.NewDecoder(r.Body)

		var data HumMessage
		err := json.NewDecoder(r.Body).Decode(&data)
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

func ReadReqBodyPOST(req *http.Request, humMess *HumMessage) {
	 body, err := ioutil.ReadAll(req.Body)
	 if err != nil {
	 	return //r1
	 }
	fmt.Println(string(body)) //DEBUG output
	//err := json.NewDecoder(req.Body).Decode(humMess)
	err = json.Unmarshal(body, humMess)
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
		fmt.Fprint(writeRespon, "400 Some errors")
	} else {

		fmt.Fprint(writeRespon, "200 Post done")
	}

	//fmt.Println("Successfully added 'The Big someNewMessage' to  table")
}

//HandlerOfMessages This func should process any messages end point
func HandlerOfMessages(w http.ResponseWriter, r *http.Request) {

	//DEBUG
///////////////	currentUserID := GetActionUserID(r)
	//fmt.Println(currentUserID)

	if r.Method == "POST" || r.Method == "OPTIONS" {
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
		fmt.Println(r.Body)
		//// https://gist.github.com/alyssaq/75d6678d00572d103106
		var inputMessage HumMessage
		inputMessage.MessageUser.IdSql = 23//currentUserID

		ReadReqBodyPOST(r, &inputMessage)

		ValidateDataFromUser(&inputMessage) //Its FAKE

		//assume data validated
		//and it is safe to put it into a Dynamo

		PutMessageToDynamo(w, &inputMessage)

	}
}
