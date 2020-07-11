package usecase

import (
	"fmt"

	"github.com/totemcaf/quongo/app/model"
	"github.com/totemcaf/quongo/app/utils"
)

type queue struct {
	queueRepository model.QueueRepository
	clock           utils.Clock
}

func NewQueueInteractor(queueRepository model.QueueRepository, clock utils.Clock) *queue {
	return &queue{queueRepository: queueRepository, clock: clock}
}

func (i *queue) FindAll(offset int, limit int) ([]*model.Queue, error) {
	return nil, utils.ErrNotImplemented
}

func (i *queue) FindByID(queueID string) (*model.Queue, error) {
	return nil, utils.ErrNotImplemented
}

func (i *queue) Complete(queue *model.Queue) *model.QueueWithStats {
	return &model.QueueWithStats{}
}

func (i *queue) Add(queue model.Queue) (*model.Queue, error) {
	return nil, utils.ErrNotImplemented
}

func (i *queue) Update(queue model.Queue) (*model.Queue, error) {
	return nil, utils.ErrNotImplemented
}

func (i *queue) FindOrCreate(queueID string) (*model.Queue, error) {
	queue, err := i.queueRepository.FindByID(queueID)

	if err == nil { // TODO Check the error is queue is not found
		return queue, nil
	}

	i.tryToCreateQueue(queueID)

	// it is read again in case other thread or instance created the queue
	return i.queueRepository.FindByID(queueID)
}

// tryToCreateQueue creates the queue and save it. It can fail if another node created the same
// queue at the same time
func (i *queue) tryToCreateQueue(queueID string) (*model.Queue, error) {
	if !model.IsQueueNameValid(queueID) {
		return nil, fmt.Errorf("Invalid queue name '%v'", queueID)
	}

	queue := model.Queue{
		Name:    queueID,
		Created: i.clock.Now(),
		VisWnd:  model.DefaultVisibilityWindow,
	}

	return i.queueRepository.Add(&queue)
}
