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
		// 允许跨域
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

// 处理所有websocket请求
func ChatRoomHandle(w http.ResponseWriter, r *http.Request) {
	conn, err := ug.Upgrade(w, r, nil)
	if err != nil {
		log.Println("here1", err)
		return
	}

	// 创建客户端
	c := &Client{
		conn: conn,
		send: make(chan Entity.Message, 128),
	}

	go c.ReadMessage()
	go c.SendMessage()
	ChatRoom.register <- c
}

// 处理所有管道任务
func (room *Room) ProcessTask() {
	log.Println("启动处理任务")
	for {
		select {
		case c := <-room.register:
			log.Println("当前有客户端进行注册")
			room.clientsPool[c] = true
		case c := <-room.unregister:
			log.Println("当前有客户端离开")
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



