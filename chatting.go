package server

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

// Room represents a single chat room
type ChatRoom struct {
	ChatRoomId    string
	publisherChan chan []byte
	subscriber    map[string](chan []byte)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 256,
}

var (
	ChatRoomId = "hum_chat_room"

	ChatRooms map[string]*ChatRoom = make(map[string]*ChatRoom)
)

func init() {
	go StartChatting(ChatRoomId)
}

func GetChatRoom(chatRoomId string) (chatRoom *ChatRoom) {

	if _, ok := ChatRooms[chatRoomId]; ok {
		// ok true
		log.Println("Join to chatting room")
	} else {
		log.Println("Chatting  Rooms aren't open ")

	}

	return ChatRooms[chatRoomId]

}

func GetSub() chan []byte {

	subChan := make(chan []byte)

	return subChan

}

func SetNewRoom(chatRoomId string) (chatRoom *ChatRoom) {

	chatRoom = new(ChatRoom)
	chatRoom.ChatRoomId = chatRoomId

	log.Println(chatRoom.ChatRoomId)
	chatRoom.publisherChan = make(chan []byte)
	chatRoom.subscriber = make(map[string](chan []byte))
	ChatRooms[chatRoomId] = chatRoom

	return

}

func StartChatting(chatRoomId string) {

	if _, ok := ChatRooms[chatRoomId]; ok {
		log.Println(chatRoomId, " has already opened")
		return
	}

	chatRoom := SetNewRoom(chatRoomId)
	log.Println(chatRoomId, " has just started")

	for {
		// waiting new messages
		msg := <-chatRoom.publisherChan

		println("new message")

		for _, sub := range chatRoom.subscriber {

			sub <- msg
			println("sending...msg to sub can")

		}

	}
}

func HandlerWs(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrading websocket ", err)
		return
	} // id_current to keep the websocket connection alive
	id_current := conn.RemoteAddr().String()

	chatRoom := GetChatRoom(ChatRoomId)

	if chatRoom == nil {
		log.Error("Chatting room isn't open")
		return
	}

	log.Infof(" Starting chatting %v\n", chatRoom)

	defer func() {
		conn.Close()
		fmt.Printf("Chatter leaving room %v\n", id_current)
		delete(chatRoom.subscriber, id_current)
	}()

	log.Infof("id : %v joined to %v chatting room\n", id_current, ChatRoomId)

	p := []byte("Hi, glad to see you here in  " + ChatRoomId)

	conn.WriteMessage(websocket.TextMessage, p)

	// Getting own chan for writting messages
	subChan := GetSub()

	chatRoom.subscriber[id_current] = subChan

	for i, sub := range chatRoom.subscriber {

		fmt.Printf("Online %v %v\n", i, sub)

	}

	go func() {
		for {

			println("waiting.......")

			if err := conn.WriteMessage(websocket.TextMessage, <-subChan); err != nil {
				log.Printf("An error happened when WriteMessage: %v", err)

				break
			}
		}

	}()

	for {

		if _, msg, err := conn.ReadMessage(); err == nil {
			chatRoom.publisherChan <- msg
		} else {
			log.Printf("An error happened when ReadMessage: %v", err)

			break
		}
	}

}