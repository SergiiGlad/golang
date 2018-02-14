package messages

import (
	"fmt"
	"go-team-room/humaws"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	//////////////////////////////////////////

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gusaul/go-dynamock"
)

var mock *dynamock.DynaMock

func init() {
	humaws.Dyna.Db, mock = dynamock.New()
}
func TestPutMessageToDynamo(t *testing.T) {
	//	lets start dynamock in action
	//	mock.ExpectGetItem().ToTable("massages").WithKeys(expectKey).WillReturns(result)
	mock.ExpectPutItem().ToTable("messages")
	someNewMessage1 := HumMessage{
		MessageID:        "1",
		MessageParentID:  "0",
		MessageTimestamp: "20180110155533001",
		MessageData: HumMessageData{
			Text:             "A lot of text and stupid smiles :)))))",
			TypeOfHumMessage: "TypeOfHumMessage-UNDEFINED FOR NOW",
			BinaryParts: []HumMessageDataBinary{HumMessageDataBinary{
				BinData: "",
				BinName: "",
			}},
		},
		MessageSocialStatus: HumMessageSocialStatus{
			Dislike: 11,
			Like:    22,
			Views:   33,
		},
		MessageUser: HumUser{
			IdSql:   777,
			NameSQL: "Vasya",
		},
	}
	respRecorder := httptest.NewRecorder()
	PutMessageToDynamo(respRecorder, &someNewMessage1)
	if status := respRecorder.Code; status != http.StatusOK {
		t.Errorf("Function PutMessageToDynamo does something wrong: got %v want %v",
			status, http.StatusOK)
	}
}

func TestGetMessageFromDynamoByUserID(t *testing.T) {
	testUserID := 777
	someNewMessage1 := HumMessage{
		MessageID:        "1",
		MessageParentID:  "0",
		MessageTimestamp: "20180110155533001",
		MessageData: HumMessageData{
			Text:             "A lot of text and stupid smiles :)))))",
			TypeOfHumMessage: "TypeOfHumMessage-UNDEFINED FOR NOW",
			BinaryParts: []HumMessageDataBinary{HumMessageDataBinary{
				BinData: "",
				BinName: "",
			}},
		},
		MessageSocialStatus: HumMessageSocialStatus{
			Dislike: 11,
			Like:    22,
			Views:   33,
		},
		MessageUser: HumUser{
			IdSql:   testUserID,
			NameSQL: "Vasya",
		},
	}
	tempMarshaledMessage, err := dynamodbattribute.MarshalMap(someNewMessage1)
	if err != nil {
		fmt.Println("Got error marshalling map:")
		fmt.Println(err.Error())
		//os.Exit(1)
	}
	//fmt.Println(tempMarshaledMessage)

	result := dynamodb.ScanOutput{
		Items: []map[string]*dynamodb.AttributeValue{tempMarshaledMessage},
	}

	mock.ExpectScan().Table("messages").WillReturns(result)
	userMessagesFromDynamo := GetMessageFromDynamoByUserID(testUserID)

	//get here without errors - it is great !!!
	//we put into mock only one message - check it
	if !reflect.DeepEqual(someNewMessage1, userMessagesFromDynamo[0]) {
		t.Errorf("Test TestGetMessageFromDynamoByUserID Fail!")
	}
}
func TestHandlerOfGetMessages(t *testing.T) {
	req, err := http.NewRequest("GET", "/messages/", nil)
	if err != nil {
		t.Fatal(err)
	}
	respRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlerOfGetMessages)

	handler.ServeHTTP(respRecorder, req)
	if status := respRecorder.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
	//
	req, err = http.NewRequest("GET", "/messages/?id=777",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	respRecorder = httptest.NewRecorder()

	handler.ServeHTTP(respRecorder, req)
	if status := respRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	//
	req, err = http.NewRequest("GET", "/messages/?id=777&maxMessages=33", nil)
	if err != nil {
		t.Fatal(err)
	}
	respRecorder = httptest.NewRecorder()

	handler.ServeHTTP(respRecorder, req)
	if status := respRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// // More test to checks for  respon body neded
	// expected := `{"user": "Vasya"}`
	// if respRecorder.Body.String() != expected {
	// 	t.Errorf("handler returned unexpected body: got %v want %v",
	// 		respRecorder.Body.String(), expected)
	// }
}

