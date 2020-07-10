package model

import (
	"time"

	"github.com/totemcaf/quongo/main/model/message"
)

// MessageRepositoryProvider ...
type MessageRepositoryProvider interface {
	ForQueue(queueID string) (MessageRepository, error)
}

// MessageRepository persist the messages of the queues
type MessageRepository interface {
	// Return the message with the provided id or nil.
	FindByID(mid string) (*Message, error)

	// Find at most 'limit' visible messages from 'offset' sorted by the original visibility time
	Find(offset int, limit int) ([]Message, error)

	/**
	 * Returns the first message that is visible and set it as taken (set its ACK value, and increment the retries)
	 * (This is provided by the repo, because it should be atomic)
	 * If no message is pending, returns nil
	 */
	Pop(quantity int) ([]*Message, error)

	// Add ads a new message to the queue
	Add(message *Message) (*Message, error)

	// Ack acknoledge the given message if ackId is the corresponding in the message
	Ack(message *Message, ackID string) error
}

// Message is the unit of data in the queue.
type Message struct {
	ID         message.MID `json:"_id"`
	Payload    string      `json:"payload"`    // Opaque data of this message
	Created    time.Time   `json:"created"`    // When the message was received and persisted
	Programmed time.Time   `json:"programmed"` // First time the message is available
	Visible    time.Time   `json:"visible"`    // When the message can be claimed for processing
	Cid        string      `json:"cid"`        // Correlation Id
	Gid        string      `json:"gid"`        // Group Id. Only one message of same group if can be in the queue
	Holder     string      `json:"holder"`     // Opaque id of the processing entity of this message
	Retries    int16       `json:"retries"`    // Number of times the message was claimed
	Ack        *string     `json:"ack"`        // Confirmation key of this message processing
	Headers    []Header
}

// Header is an additional metadata associated with the message
// It is used to transport additional semantic to the message interpreted by the consumer
type Header struct {
	name  string
	value string
}
