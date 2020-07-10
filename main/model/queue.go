package model

import (
	"regexp"
	"time"
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
	Push    bool          `json:"push"`
}

var queueNamePattern = regexp.MustCompile(`^[-a-zA-Z0-9_]+$`)

// IsQueueNameValid Report if name is valid queue name
func (queue Queue) IsQueueNameValid(name string) bool {
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
