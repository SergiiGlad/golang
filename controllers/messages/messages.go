package messages

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Create structs to hold info about new item of HumMessage

type HumUser struct {
	IdSql   int    `json:"id_sql"`
	NameSql string `json:"name_sql"`
}

type HumMessageDataBinary struct {
	BinData string `json:"bin_data"`
	BinName string `json:"bin_name"`
}

type HumMessageData struct {
	Text             string                 `json:"text"`
	TypeOfHumMessage string                 `json:"type"`
	BinaryParts      []HumMessageDataBinary `json:"binary_parts"`
}
type HumMessageSocialStatus struct {
	Dislike int
	Like    int
	Views   int
}

type HumMessage struct {
	MessageId           string                 `json:"message_id"`
	MessagechatRoomId   string                 `json:"message_chat_room_id`
	MessageData         HumMessageData         `json:"message_data"`
	MessageParentId     string                 `json:"message_parent_id"`
	MessageSocialStatus HumMessageSocialStatus `json:"message_social_status"`
	MessageTimestamp    string                 `json:"message_timestamp"`
	MessageUser         HumUser                `json:"message_user"`
	/////////END
}

///////////////

type test_struct struct {
	//Test string
	ChatRoomId string
	Message    string
}

//GetActionUserId - return ID of user who make this action
func GetActionUserId(r *http.Request) string {
	// TEMPORARY stub

	if r.Method == "GET" {
		return "12" //later We will get it REALY

	} else if r.Method == "POST" {
		return "23"
	} else {
		return "34"
	}

}

//ValideteDataFromUser very important func
//Do NOT trust any data from User!!!
//VALIDATE EVERUTHING
func ValideteDataFromUser(m *HumMessage) {
	//work this out LATER
	//for now it is stub only
	//
	//do nothing
	//go home
}

//HandlerOfMessages This func should process any messages end point
func HandlerOfMessages(w http.ResponseWriter, r *http.Request) {
	currentUserID := GetActionUserId(r)
	fmt.Println(currentUserID)
	//UserIDfromGET := r.Form.Get("id") // id will be "" if parameter is not set
	if r.Method == "GET" {
		r.ParseForm() // Parses the request body

		//Assume it is an GET
		//Shuold return an  last messages
		//for user ID=currentUserID

	} else if r.Method == "POST" || r.Method == "OPTIONS" {
		//CORS!!! "C:\Program Files (x86)\Google\Chrome\Application\chrome.exe" --disable-web-security --user-data-dir="D:/Chrome"
		//Assume it is an POST
		//Shuold  expect a user messagein form of:
		// {
		// 	"id": 0,
		// 	"sendDate": "string",
		// 	"isNew": true,
		// 	"from": 0,
		// 	"body": {
		// 	  "chatRoomId": 0,
		// 	  "value": {
		// 		"message": "string",
		// 		"fileDescription": "string",
		// 		"file": "string"
		// 	  }
		// 	}
		//   }

		/////////////////////
		//fmt.Println(r.Body)
		fmt.Println("===========")
		//
		//// https://gist.github.com/alyssaq/75d6678d00572d103106
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}
		//var postDataFromRequest string
		//postDataFromRequest = postDataFromRequest + string(body)
		fmt.Println(string(body))
		var inputMessage HumMessage
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&inputMessage)

		if err != nil {
			panic(err)
		}
		ValideteDataFromUser(&inputMessage)
		fmt.Println(inputMessage.MessageData)

		////////////////////////////////////
		//fmt.Fprint(w, "POST done")

	}
}
