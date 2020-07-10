package memory

import (
	"errors"
	"log"
	"sync"

	"github.com/golang-collections/go-datastructures/queue"
	"github.com/totemcaf/quongo/main/model"
)

type item struct {
	message *model.Message
}

// MessageRepository ...
type MessageRepository struct {
	queueID         string
	queueRepository model.QueueRepository
	lock            sync.RWMutex
	messages        queue.PriorityQueue
}

// NewMessageRepository ...
func NewMessageRepository() *MessageRepository {
	return &MessageRepository{}
}

// Find finds at most 'limit' visible messages from 'offset' sorted by the original visible time
func (r *MessageRepository) Find(offset int, limit int) ([]model.Message, error) {
	result := make([]model.Message, 0, 1)

	return result, nil
}

// FindByID returns the message with the provided id or nil.
func (r *MessageRepository) FindByID(mid string) (*model.Message, error) {
	var result model.Message

	return &result, nil
}

// Pop returns the first message that is visible and set it as taken (set its ACK value, and increment the retries)
// If no message is pending, returns nil
func (r *MessageRepository) Pop(quantity int) ([]*model.Message, error) {
	result := make([]*model.Message, 0)
	return result, nil
}

// Add a message into the queue
func (r *MessageRepository) Add(message *model.Message) (*model.Message, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.messages.Put(&item{message})

	return message, nil
}

// Ack ...
func (r *MessageRepository) Ack(message *model.Message, ackID string) error {
	return errors.New("Not implemented")
}

func (i item) Compare(other queue.Item) int {
	item2, ok := other.(item)

	if !ok {
		log.Fatal("Receive incorrect type")
	}

	if i.message.Programmed.Before(item2.message.Programmed) {
		return -1
	}

	if i.message.Programmed.After(item2.message.Programmed) {
		return 1
	}

	return 0
}
