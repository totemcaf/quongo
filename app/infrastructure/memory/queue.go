package memory

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/totemcaf/quongo/app/model"
	"github.com/totemcaf/quongo/app/utils"
)

// Repository ...
type Repository struct {
	queues map[string]*model.Queue
	clock  utils.Clock
	lock   sync.RWMutex
}

// NewQueueRepository ...
func NewQueueRepository(clock utils.Clock) *Repository {
	return &Repository{queues: make(map[string]*model.Queue, 0), clock: clock}
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

	r.lock.RLock()
	defer r.lock.RUnlock()

	queue, found := r.queues[id]

	if found {
		log.Printf("Found: %v", queue.Name)
		return queue, nil
	}

	log.Printf("Not found")
	return nil, fmt.Errorf("Not found queue '%v'", id)
}

// Add ...
func (r *Repository) Add(queue *model.Queue) (*model.Queue, error) {
	log.Printf("Adding queue %v", queue.Name)

	r.lock.Lock()
	defer r.lock.Unlock()

	_, ok := r.queues[queue.Name]

	if ok {
		return nil, fmt.Errorf("Duplicate queue '%v'", queue.Name)
	}

	r.queues[queue.Name] = queue.
		WithRepository(NewMessageRepository(queue.Name, r.clock)).
		WithClock(r.clock)

	return queue, nil
}

// Update ...
func (r *Repository) Update(queue *model.Queue) (*model.Queue, error) {
	return queue, nil // TODO do it
}

// ForQueue ...
func (r *Repository) ForQueue(queueID string) (model.MessageRepository, error) {
	return nil, nil
}
