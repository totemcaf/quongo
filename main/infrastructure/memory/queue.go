package memory

import (
	"errors"
	"log"
	"sync"

	"github.com/totemcaf/quongo/main/model"
)

type queueEntry struct {
	metadata    model.Queue
	messageRepo model.MessageRepository
}

// Repository ...
type Repository struct {
	queues map[string]queueEntry
	lock   sync.RWMutex
}

// NewQueueRepository ...
func NewQueueRepository() *Repository {
	return &Repository{}
}

// FindAll Reports at most 'size' queues starting from 'skip'
func (r *Repository) FindAll(skip int, size int) ([]*model.Queue, error) { // TODO Should be *model.Queue ??
	if skip < 0 {
		return nil, errors.New("Invalid skip, it should be >= 0")
	}

	if size < 0 {
		return nil, errors.New("Invalid size, it shoud be >= 0")
	}

	if size == 0 {
		return []*model.Queue{}, nil
	}

	result := make([]*model.Queue, size) // TODO Should be "*" or struct directly?

	return result, nil
}

// Complete Adds the queue statistics to the provided queue
func (r *Repository) Complete(queue *model.Queue) *model.QueueWithStats {
	hidden, visible := 0, 0

	return model.NewQueueWithStat(queue, visible+hidden, hidden, visible)
}

// FindByID Returns a queue with the provided id (name)
func (r *Repository) FindByID(id string) (*model.Queue, error) {
	log.Printf("Finding queue by id: %v", id)
	result := model.Queue{}

	err := errors.New("Not implemented")

	if err == nil {
		log.Printf("Found: %v", result.Name)
		return &result, nil
	}

	log.Printf("Not found")
	return nil, err
}

// Add ...
func (r *Repository) Add(queue *model.Queue) (*model.Queue, error) {
	log.Printf("Inserting queue %v", queue.Name)
	err := errors.New("Not implemented")
	if err != nil {
		log.Printf("Error:      %v", err)
	}
	return queue, err
}

// Update ...
func (r *Repository) Update(queue *model.Queue) (*model.Queue, error) {
	return queue, nil // TODO do it
}

// ForQueue ...
func (r *Repository) ForQueue(queueID string) (model.MessageRepository, error) {
	if queueID == "" {
		return nil, errors.New("Invalid queue id, it is empty")
	}

	queue, found := r.findExistentQueue(queueID)

	if found {
		return queue.messageRepo, nil
	}

	return r.createQueueWithDefaults(queueID).messageRepo, nil
}

func (r *Repository) findExistentQueue(queueID string) (queueEntry, bool) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	entry, err := r.queues[queueID]
	return entry, err
}

func (r *Repository) createQueueWithDefaults(queueID string) queueEntry {
	r.lock.Lock()

	// Check again in case another coroutine won and create it before me
	entry, found := r.queues[queueID]

	if found {
		return entry
	}

	entry = queueEntry{
		metadata:    model.Queue{},
		messageRepo: NewMessageRepository(),
	}

	return entry
}
