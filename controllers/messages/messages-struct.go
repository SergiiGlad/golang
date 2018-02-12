package messages

// Create structs to hold info about new item of HumMessage

//HumUser simple structure to hold Hum user mame and ID from sql DB
type HumUser struct {
	IdSql   int    `json:"id_sql"`
	NameSQL string `json:"name_sql"`
}

//HumMessageDataBinary - struct to hold data about binary part of message
type HumMessageDataBinary struct {
	BinData string `json:"bin_data"`
	BinName string `json:"bin_name"`
}

//HumMessageData simple structure to hold HumMessage date like type og message and text
type HumMessageData struct {
	Text             string                 `json:"text"`
	TypeOfHumMessage string                 `json:"type"`
	BinaryParts      []HumMessageDataBinary `json:"binary_parts"`
}

//HumMessageSocialStatus - simple structure to hold HumMessage (un)likes\views
type HumMessageSocialStatus struct {
	Dislike int
	Like    int
	Views   int
}

//HumMessage complete structure of HumMessage
type HumMessage struct {
	MessageID           string                 `json:"message_id"`
	MessageChatRoomID   string                 `json:"message_chat_room_id"`
	MessageData         HumMessageData         `json:"message_data"`
	MessageParentID     string                 `json:"message_parent_id"`
	MessageSocialStatus HumMessageSocialStatus `json:"message_social_status"`
	MessageTimestamp    string                 `json:"message_timestamp"`
	MessageUser         HumUser                `json:"message_user"`
	/////////END
}

//HumChatRoom - complete struct of HumChatRoom
type HumChatRoom struct {
	ChatAdminUserID  []HumUser `json:"chat_admin_user_id"`
	ChatCreationDate string    `json:"chat_creation_date"`
	ChatID           string    `json:"chat_id"`
	ChatName         string    `json:"chat_name"`
	ChatStatus       string    `json:"chat_status"`
	ChatTitile       string    `json:"chat_title"`
	ChatUsersList    []HumUser `json:"chat_users_list"`
}
