package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Chat struct {
	ID        primitive.ObjectID  `json:"_id,omitempty" bson:"_id,omitempty"`
	Message   string              `json:"message,omitempty"`
	Sender    primitive.ObjectID  `json:"sender,omitempty" bson:"sender,omitempty"`
	Type      string              `json:"type,omitempty"`
	CreatedAt primitive.Timestamp `json:"createdAt"`
	UpdatedAt primitive.Timestamp `json:"updatedAt"`
}
