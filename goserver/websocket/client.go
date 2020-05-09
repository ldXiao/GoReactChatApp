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
		var chat_info map[string]string
		// json.NewDecoder(p.(io.Reader)).Decode(&chat)
		// message := models.Chat{
		// 	Type:    chat.Type,
		// 	Message: ,
		// }
		err = json.Unmarshal(p, &chat_info)
		if err != nil {
			log.Fatal("failed")
		}
		fmt.Println(chat_info)
		objid, err := primitive.ObjectIDFromHex(chat_info["userId"])
		var chat = models.Chat{
			Type:    chat_info["type"],
			Message: chat_info["chatMessage"],
			Sender:  objid,
		}

		c.Pool.Broadcast <- chat
		fmt.Printf("Message Received: %+v\n", chat)

	}
}
