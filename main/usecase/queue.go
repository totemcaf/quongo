package usecase

import (
	"errors"
	"github.com/totemcaf/quongo/main/model"
	"log"
	"time"
)

type QueueInteractor struct {
	repository model.QueueRepository
}

func NewQueueInteractor(repository model.QueueRepository) *QueueInteractor {
	return &QueueInteractor{repository: repository}
}

func (i *QueueInteractor) FindAll(offset int, limit int) ([]*model.Queue, error) {
	return i.repository.FindAll(offset, limit)
}

func (i *QueueInteractor) FindById(queueId string) (*model.Queue, error) {
	return i.repository.FindById(queueId)
}

func (i *QueueInteractor) Complete(queue *model.Queue) model.QueueWithStats {
	return model.QueueWithStats{} // TODO
}

func (i *QueueInteractor) IsQueueNameValid(queueId string) bool {
	return false // TODO
}

func (i *QueueInteractor) Add(queue model.Queue) (*model.Queue, error) {
	log.Printf("Adding %v", queue.Name)

	if !queue.IsQueueNameValid(queue.Name) {
		return nil, errors.New("queue name is invalid")
	}

	queue.Created = time.Now()

	if queue.VisWnd <= 0 {
		queue.VisWnd = model.DefaultVisibilityWindow // TODO Queue
	}

	return i.repository.Add(&queue)
}

func (i *QueueInteractor) Update(queue model.Queue) (*model.Queue, error) {
	return nil, nil // TODO
}
