package model

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type MessageRepository interface {
	/**
	 * Return the message with the provided id or nil.
	 */
	FindById(mid string) (*Message, error)

	/**
	 * Find at most 'limit' visible messages from 'offset' sorted by the original visibility time
	 */
	Find(offset int, limit int) ([]Message, error)

	/**
	 * Returns the first message that is visible and set it as taken (set its ACK value, and increment the retries)
	 * (This is provided by the repo, because it should be atomic)
	 * If no message is pending, returns nil
	 */
	Pop(delay time.Duration) (*Message, error)
	Add(message *Message) (*Message, error)
}

type Message struct {
	Id         bson.ObjectId  `json:"_id"`
	Payload    string         `json:"payload,string"` // Opaque data of this message
	Created    time.Time      `json:"created"`        // When the message was received and persisted
	Programmed time.Time      `json:"programmed"`     // First time the message is available
	Visible    time.Time      `json:"visible"`        // When the message can be claimed for processing
	Cid        string         `json:"cid"`            // Correlation Id
	Gid        string         `json:"gif"`            // Group Id. Only one message of same group if can be in the queue
	Holder     string         `json:"holder"`         // Opaque id of the processing entity of this message
	Retries    int16          `json:"retries"`        // Number of times the message was claimed
	Ack        *bson.ObjectId `json:"ack"`            // Confirmation key of this message processing
}
