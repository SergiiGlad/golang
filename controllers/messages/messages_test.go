package messages

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gusaul/go-dynamock"
)

func TestPutMessageToDynamo(t *testing.T) {
	var mock *dynamock.DynaMock
	Dyna.Db, mock = dynamock.New()
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

	//	lets start dynamock in action
	mock.ExpectGetItem().ToTable("massages").WithKeys(expectKey).WillReturns(result)

}

func TestGetMessageFromDynamoByUserID(t *testing.T) {
	testUserID := 777
	// someNewMessage1 := HumMessage{
	// 	MessageId:        "1",
	// 	MessageParentId:  "0",
	// 	MessageTimestamp: "20180110155533001",
	// 	MessageData: HumMessageData{
	// 		Text:             "A lot of text and stupid smiles :)))))",
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
	// 		IdSql:   777,
	// 		NameSql: "Vasya",
	// 	},
	// }
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

	GetMessageFromDynamoByUserID(testUserID)

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
