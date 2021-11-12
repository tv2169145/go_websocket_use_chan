package roomHandler

import (
	"github.com/gorilla/websocket"
	"github.com/tv2169145/go_websocket_use_chan/Entity"
	"log"
	"time"
)

// 客户端
type Client struct {
	conn *websocket.Conn
	send chan Entity.Message
}

// 接收消息
func (c *Client) ReadMessage() {
	preMessageTime := int64(0)
	for {
		message := &Entity.Message{}
		if err := c.conn.ReadJSON(message); err != nil {
			c.conn.Close()
			ChatRoom.unregister <- c
			return
		}

		// 限制用户发送消息频率，每1秒只能发送一条消息
		curMessageTime := time.Now().Unix()
		if curMessageTime-preMessageTime < 0 {
			log.Println("1秒内不可重發")
			continue
		}
		preMessageTime = curMessageTime
		ChatRoom.send <- *message
	}
}

// 發送消息
func (c *Client) SendMessage() {
	for {
		m := <-c.send
		if err := c.conn.WriteJSON(m); err != nil {
			c.conn.Close()
			ChatRoom.unregister <- c
			return
		}
	}
}
