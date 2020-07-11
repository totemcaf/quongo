package model

import (
	"regexp"
	"time"

	"github.com/totemcaf/quongo/app/utils"
)

// MaxQueueNameLen ...
const MaxQueueNameLen = 32

// DefaultVisibilityWindow ...
const DefaultVisibilityWindow = 30 * time.Second

// QueueRepository ...
type QueueRepository interface {
	// Report at most 'size' queues starting from 'skip'
	FindAll(skip int, size int) ([]*Queue, error)

	// Return the message with the provided id or nil.
	FindByID(queueID string) (*Queue, error)
	Add(queue *Queue) (*Queue, error)
	Complete(queue *Queue) *QueueWithStats
}

// Queue ...
type Queue struct {
	Name    string        `json:"name" bson:"_id"`
	Created time.Time     `json:"created"`
	VisWnd  time.Duration `json:"visibilityWindow" bson:"visibilityWindow"`

	// Local variables
	messageRepository MessageRepository
	clock             utils.Clock
}

var queueNamePattern = regexp.MustCompile(`^[-a-zA-Z0-9_]+$`)

// IsQueueNameValid Report if name is valid queue name
func IsQueueNameValid(name string) bool {
	return len(name) > 0 && len(name) <= MaxQueueNameLen && queueNamePattern.MatchString(name)
}

// QueueWithStats ...
type QueueWithStats struct {
	*Queue `json:"queue"`
	Stats  QueueStats `json:"stats"`
}

// QueueStats ...
type QueueStats struct {
	Total     int `json:"total"`
	Hidden    int `json:"hidden"`
	InProcess int `json:"inProcess"`
}

// NewQueueWithStat ...
func NewQueueWithStat(queue *Queue, total int, hidden int, inProcess int) *QueueWithStats {
	return &QueueWithStats{queue, QueueStats{total, hidden, inProcess}}
}

// Add a message to this queue
func (q *Queue) Add(message *Message) (*Message, error) {
	return q.messageRepository.Add(message)
}

// Pop get a quantity of available messages and lock them
func (q *Queue) Pop(quantity int) ([]*Message, error) {
	messages, err := q.messageRepository.PopAvailable(quantity)

	if err != nil {
		return nil, err
	}

	// WARNING this is not atomic
	for _, msg := range messages {
		msg.DelayTo(q.clock.Now().Add(q.VisWnd)).Lock() // TODO Impove lock all at once
		q.messageRepository.Update(msg)                 // TODO Improve with UpdateAll
	}

	return messages, nil
}

// WithRepository augments the queue with the repository for it
func (q *Queue) WithRepository(messageRepository MessageRepository) *Queue {
	return &Queue{
		q.Name,
		q.Created,
		q.VisWnd,
		messageRepository,
		q.clock,
	}
}

// WithClock augments the queue with the clock to use
func (q *Queue) WithClock(clock utils.Clock) *Queue {
	return &Queue{
		q.Name,
		q.Created,
		q.VisWnd,
		q.messageRepository,
		clock,
	}
}
