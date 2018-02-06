package messages

// Create structs to hold info about new item of HumMessage

//
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
	MessageChatRoomId   string                 `json:"message_chat_room_id"`
	MessageData         HumMessageData         `json:"message_data"`
	MessageParentId     string                 `json:"message_parent_id"`
	MessageSocialStatus HumMessageSocialStatus `json:"message_social_status"`
	MessageTimestamp    string                 `json:"message_timestamp"`
	MessageUser         HumUser                `json:"message_user"`
	/////////END
}
