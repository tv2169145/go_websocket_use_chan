package roomHandler

import (
	"github.com/gorilla/websocket"
	"github.com/tv2169145/go_websocket_use_chan/Entity"
	"log"
	"net/http"
)

var (
	ChatRoom *Room
	ug       = websocket.Upgrader{
		// 允許跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// 聊天室配置
type Room struct {
	register    chan *Client
	unregister  chan *Client
	clientsPool map[*Client]bool
	send        chan Entity.Message
}

func NewRoom() {
	ChatRoom = &Room{
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		clientsPool: map[*Client]bool{},
		send:        make(chan Entity.Message),
	}
}

func Broadcast(w http.ResponseWriter, r *http.Request) {
	ChatRoom.send <- Entity.Message{Token:"Reid", Content: "Hello"}
	w.Write([]byte("ok"))
}

// 處理所有websocket請求
func ChatRoomHandle(w http.ResponseWriter, r *http.Request) {
	conn, err := ug.Upgrade(w, r, nil)
	if err != nil {
		log.Println("here1", err)
		return
	}

	// 建立client
	c := &Client{
		conn: conn,
		send: make(chan Entity.Message, 50000),
	}

	go c.ReadMessage()
	go c.SendMessage()
	ChatRoom.register <- c
}

// 處理所有管道任务
func (room *Room) ProcessTask() {
	log.Println("do process")
	for {
		select {
		case c := <-room.register:
			log.Println("on connection")
			room.clientsPool[c] = true
		case c := <-room.unregister:
			log.Println("disconnection")
			if room.clientsPool[c] {
				close(c.send)
				delete(room.clientsPool, c)
			}
		case m := <-room.send:
			for c := range room.clientsPool {
				select {
				case c.send <- m:
				default:
					break
				}
			}
		}
	}
}

// 启动聊天室
//func StartChatRoom() {
//	log.Println("聊天室启动....")
//	mux := pat.New()
//
//	ChatRoom = &Room{
//		register:    make(chan *Client),
//		unregister:  make(chan *Client),
//		clientsPool: map[*Client]bool{},
//		send:        make(chan Entity.Message),
//	}
//	mux.Get("/", http.HandlerFunc(homeHandler.Home))
//	mux.Get("/ws", http.HandlerFunc(chatRoomHandle))
//
//	go chatRoom.ProcessTask()
//	_ = http.ListenAndServe(":8080", mux)
//	//if err := http.ListenAndServe(":8081", nil); err != nil {
//	//	panic(err)
//	//}
//}



