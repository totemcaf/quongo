package memory

import (
	"errors"
	"sort"
	"sync"

	"github.com/totemcaf/quongo/app/model"
	"github.com/totemcaf/quongo/app/utils"
)

// MessageRepository ...
type MessageRepository struct {
	queueID         string
	queueRepository model.QueueRepository
	lock            sync.RWMutex
	clock           utils.Clock
	messages        []*model.Message
}

// NewMessageRepository ...
func NewMessageRepository(queueID string, clock utils.Clock) *MessageRepository {
	return &MessageRepository{
		queueID: queueID,
		clock:   clock,
	}
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
func (r *MessageRepository) PopAvailable(quantity int) ([]*model.Message, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	result := make([]*model.Message, 0, quantity)

	now := r.clock.Now()

	for _, msg := range r.messages {
		if len(result) >= quantity || now.Before(msg.Visible) {
			// No more messages
			return result, nil
		}
		result = append(result, msg)
	}

	return result, nil
}

// Add a message into the queue
func (r *MessageRepository) Add(message *model.Message) (*model.Message, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	ms := r.messages

	if ms == nil || len(ms) == 0 {
		ms = make([]*model.Message, 0, 100)
		ms = append(ms, message)
		r.messages = ms
		return message, nil
	}

	index := findSortedPosition(ms, message)

	if index > len(ms) {
		r.messages = append(ms, message)
	} else {
		ms = append(ms, nil)

		copy(ms[index+1:], ms[index:])
		ms[index] = message
		r.messages = ms
	}

	return message, nil
}

func findPosition(items []*model.Message, msg *model.Message) int {
	if items != nil {
		mid := msg.ID

		for i, m := range items {
			if m.ID == mid {
				return i
			}
		}
	}

	return -1
}

func findSortedPosition(items []*model.Message, msg *model.Message) int {
	return sort.Search(len(items), func(i int) bool {
		return compare(items[i], msg) >= 0
	})
}

// Update store the new values of the message keeping the sort order
func (r *MessageRepository) Update(message *model.Message) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	ms := r.messages
	if ms != nil {
		oldIdx := findPosition(ms, message)

		if oldIdx >= 0 {
			newIdx := findSortedPosition(ms, message)

			if oldIdx < newIdx {
				// Should move elements to the left
				copy(ms[oldIdx:newIdx-1], ms[oldIdx+1:newIdx])
			} else if newIdx < oldIdx {
				// Should move elements to the right
				copy(ms[newIdx+1:oldIdx], ms[newIdx:oldIdx-1])
			}

			ms[newIdx] = message

			return nil
		}
	}

	return errors.New("Message not found")
}

// Ack ...
func (r *MessageRepository) Ack(message *model.Message, ackID string) error {
	return errors.New("Not implemented")
}

func compare(m1 *model.Message, m2 *model.Message) int {
	if m1.Visible.Before(m2.Visible) {
		return -1
	}

	if m1.Visible.After(m2.Visible) {
		return 1
	}

	return 0
}
