package main

import "time"

type Message struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	Content   string    `json:"content" bson:"content"`
	Receiver  string    `json:"receiver" bson:"receiver"`
	Sender    string    `json:"sender" bson:"sender"`
	Room      string    `json:"room" bson:"room"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}
