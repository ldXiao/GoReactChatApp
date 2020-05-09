package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/ldXiao/GoReactChatApp/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
	mu   sync.Mutex
}

// type Message struct {
// 	Type int    `json:"type"`
// 	Body string `json:"body"`
// }

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		var chat_map map[string]string
		err = json.Unmarshal(p, &chat_map)
		if err != nil {
			log.Fatal("failed")
		}
		fmt.Println(chat_map)
		objid, err := primitive.ObjectIDFromHex(chat_map["userId"])
		var chat = models.Chat{
			Type:    chat_map["type"],
			Message: chat_map["chatMessage"],
			Sender:  objid,
		}
		chat.Save()
		chat_info := chat.GetChatInfo()
		c.Pool.Broadcast <- chat_info
		fmt.Printf("Message Received: %+v\n", chat_info)

	}
}
