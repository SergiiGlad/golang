// chat example slightly to represent real world chatrooms

package server

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"time"

	"go-team-room/controllers/messages"
)

// Room represents a single chat room
type ChatRoom struct {
	ChatRoomId    string
	publisherChan chan []byte
	subscriber    map[string](chan []byte)
}

func init() {

	// need to start chatting room
	go StartChatting(ChatRoomId)
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size, here 1 kB
	maxMessageSize = 1024 * 1024
)

var (
	ChatRoomId                      = "hum_chat_room"
	ChatRooms  map[string]*ChatRoom = make(map[string]*ChatRoom)
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 256,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Getting char room by chatRoomId from map ChatRooms

func GetChatRoom(chatRoomId string) (chatRoom *ChatRoom) {

	if _, ok := ChatRooms[chatRoomId]; ok {
		// ok true
		log.Println("New chatter join to room ", chatRoomId)
	} else {
		log.Println("No Chatting Rooms have been started")

	}

	return ChatRooms[chatRoomId]

}

// Initial chatting room

func SetNewRoom(chatRoomId string) (chatRoom *ChatRoom) {

	chatRoom = new(ChatRoom)
	chatRoom.ChatRoomId = chatRoomId

	fmt.Println(chatRoom.ChatRoomId)
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

		log.Println("a new message received")

		for _, sub := range chatRoom.subscriber {

			log.Println("sending messages to all subcribers")
			sub <- msg
		}

	}
}

// Greate new chan for a subcriber

func GetSub() chan []byte {

	subChan := make(chan []byte)

	return subChan

}

func SayWelcome(c *websocket.Conn, privet *messages.HumMessage) error {
	return c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%v", *privet)))

}

func PrintChatters(chatRoom *ChatRoom) {

	for i, sub := range chatRoom.subscriber {

		log.Printf("Online %v %v\n", i, sub)

	}
}

func GetUserHumMessage(privet string) *messages.HumMessage {
	return &messages.HumMessage{
		MessageID: "1",
	}
}

func HandlerWs(w http.ResponseWriter, r *http.Request) {

	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrader failed with error : ", err)
		return
	} else {
		log.Println("Handler represented a websocket connection")
	}

	// getting chatroom
	chatRoom := GetChatRoom(ChatRoomId)

	if chatRoom == nil {
		log.Println("Chatting room isn't open")
		return
	}

	// id_current for each chatters. Storing in map chatRoom.subscriber
	id_current := conn.RemoteAddr().String() + r.UserAgent()

	// Getting own chan for writting messages and add chan to map of all subscribers
	subChan := GetSub()
	chatRoom.subscriber[id_current] = subChan

	defer func() {

		fmt.Printf("chatter leaving room %v\n", id_current)
		conn.WriteControl(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, "socket close"),
			time.Now().Add(time.Second/10))
		conn.Close()
		delete(chatRoom.subscriber, id_current)
	}()

	// first HumMessage
	chatter_name := "Chatter"
	firstMsg := GetUserHumMessage("Hi, glad to see you here in  " + ChatRoomId)
	firstMsg.MessageUser.NameSQL = chatter_name
	firstMsg.MessageChatRoomID = ChatRoomId

	// logging all chatters
	PrintChatters(chatRoom)

	log.Printf(" User %v is starting chatting in %v\n", chatter_name, chatRoom)

	SayWelcome(conn, firstMsg)

	go func() {
		for {

			// waiting writting new messages...

			if err := conn.WriteMessage(websocket.TextMessage, <-subChan); err != nil {
				log.Printf("An error happened when WriteMessage: %v", err)
				time.Sleep(time.Second)
				break
			}
		}

	}()

	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {

		if messageType, msg, err := conn.ReadMessage(); err != nil {

			if websocket.IsUnexpectedCloseError(err) {
				log.Printf("Unexpected Close WebSocket Conection: %v", err)

			} else {

				log.Printf("An error happened when Read Message: %v", err)
			}
			time.Sleep(time.Second)
			break

		} else {

			switch messageType {
			case websocket.CloseNormalClosure:
				log.Printf("Web Socket Normal Close: %v")

			case websocket.TextMessage:
				chatRoom.publisherChan <- msg
			}

		}

	}

}
