package messages

import (
	"net/http"
	"net/http/httptest"
	"testing"
//////////////////////////////////////////
	"github.com/gusaul/go-dynamock"

)

var mock *dynamock.DynaMock

func init() {
	Dyna.Db, mock = dynamock.New()
}
func TestPutMessageToDynamo(t *testing.T) {

	// expectKey := map[string]*dynamodb.AttributeValue{
	// 	"id": {
	// 		N: aws.String("23"),
	// 	},
	// }
	// expectedResult := aws.String("Vasya")
	// result := dynamodb.GetItemOutput{
	// 	Item: map[string]*dynamodb.AttributeValue{
	// 		"name": {
	// 			S: expectedResult,
	// 		},
	// 	},
	// }

	//	lets start dynamock in action
	//	mock.ExpectGetItem().ToTable("massages").WithKeys(expectKey).WillReturns(result)
	//mock.ExpectGetItem().ToTable("massages").WithKeys().WillReturns()
	mock.ExpectPutItem().ToTable("massages")
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

/*
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
			IdSql:   777,
			NameSQL: "Vasya",
		},
	}
	// someNewMessage2 := HumMessage{
	// 	MessageId:        "2",
	// 	MessageParentId:  "0",
	// 	MessageTimestamp: "20180110155533011",
	// 	MessageData: HumMessageData{
	// 		Text:             "MORE of text and stupid smiles :)))))",
	// 		TypeOfHumMessage: "TypeOfHumMessage-UNDEFINED FOR NOW",
	// 		BinaryParts: []HumMessageDataBinary{HumMessageDataBinary{
	// 			BinData: "",
	// 			BinName: "",
	// 		}},
	// 	},
	// 	MessageSocialStatus: HumMessageSocialStatus{
	// 		Dislike: 11,
	// 		Like:    22,
	// 		Views:   33,
	// 	},
	// 	MessageUser: HumUser{
	// 		IdSql:   666,
	// 		NameSql: "Petya",
	// 	},
	// }
	tempMarshaledMessage, _ := json.Marshal(someNewMessage1)
	fmt.Println(tempMarshaledMessage)
	expectKey := map[string]*dynamodb.AttributeValue{
		"MessageUser.IdSql": {
			N: aws.String("777"),
		},
	}
	//expectedResult := aws.String("jaka")
	expectedResult := `{"message_id":"20180110155343158","message_chat_room_id":"995","message_data":{"text":"8 A lot of text and stupid smiles :)))))","type":"TypeOfHumMessage-UNDEFINED FOR NOW","binary_parts":[{"bin_data":"","bin_name":""}]},"message_parent_id":"","message_social_status":{"Dislike":15,"Like":262,"Views":373},"message_timestamp":"20180110155533111","message_user":{"id_sql":33,"name_sql":"Vasya"}}`

	result := dynamodb.GetItemOutput{
		Item: map[string]*dynamodb.AttributeValue{
			"name": {
				S: expectedResult,
			},
		},
	}
	//lets start dynamock in action
	mock.ExpectPutItem().ToTable("messages").WithKeys(expectKey).WillReturns(result)

	GetMessageFromDynamoByUserID(testUserID)

}

*/
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
