package models

import (
	"context"
	"fmt"
	"log"

	"github.com/ldXiao/GoReactChatApp/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatInterface interface {
	isChat()
}

type Chat struct {
	ID        primitive.ObjectID  `json:"_id,omitempty" bson:"_id,omitempty"`
	Message   string              `json:"message,omitempty"`
	Sender    primitive.ObjectID  `json:"sender,omitempty" bson:"sender,omitempty"`
	Type      string              `json:"type,omitempty"`
	CreatedAt primitive.Timestamp `json:"createdAt"`
	UpdatedAt primitive.Timestamp `json:"updatedAt"`
}

func (c Chat) isChat() {}

type Sender struct {
	Name  string `json:"name,omitempty"`
	Image string `json:"image, omitempty"`
}

type Chat_info struct {
	ID        primitive.ObjectID  `json:"_id,omitempty" bson:"_id,omitempty"`
	Message   string              `json:"message,omitempty"`
	Sender    Sender              `json:"sender,omitempty" bson:"sender,omitempty"`
	Type      string              `json:"type,omitempty"`
	CreatedAt primitive.Timestamp `json:"createdAt"`
	UpdatedAt primitive.Timestamp `json:"updatedAt"`
}

func (c Chat_info) isChat() {}

func (chat *Chat) Save() {
	// id and timestamps are not initialized at this step
	doc, err := middleware.ChatsCollection.InsertOne(context.Background(), chat)
	if err != nil {
		fmt.Println(1)
		log.Println(err)
	}
	chatid := doc.InsertedID.(primitive.ObjectID)
	singres := middleware.ChatsCollection.FindOne(context.TODO(), bson.D{{"_id", chatid}})
	if singres.Err() != nil {
		fmt.Println(3)
		log.Println(singres.Err())
	}
	singres.Decode(chat)
}

func (chat *Chat) GetChatInfo() Chat_info {
	var user User
	userfind := middleware.UsersCollection.FindOne(context.TODO(), bson.D{{"_id", chat.Sender}})

	if userfind.Err() != nil {
		fmt.Println(2)
		log.Println(userfind.Err())
	}
	userfind.Decode(&user)

	var chatinfo = Chat_info{
		ID:        chat.ID,
		Message:   chat.Message,
		CreatedAt: chat.CreatedAt,
		UpdatedAt: chat.UpdatedAt,
		Sender: Sender{
			Name:  user.Name,
			Image: user.Image,
		},
		Type: chat.Type,
	}
	return chatinfo
}
