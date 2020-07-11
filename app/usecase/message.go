package usecase

import (
	"errors"

	"github.com/totemcaf/quongo/app/model"
	"github.com/totemcaf/quongo/app/utils"
)

// ErrNoMessageAvailable marks there are no messages on the queue
var ErrNoMessageAvailable = errors.New("No messages available")

type queueService interface {
	FindOrCreate(queueID string) (*model.Queue, error)
}

type messageInteractor struct {
	queueServie queueService
}

func NewMessageInteractor(queueService queueService) *messageInteractor {
	return &messageInteractor{
		queueServie: queueService,
	}
}

func (i *messageInteractor) FindAll(queueID string, offset int, limit int) ([]model.Message, error) {
	return nil, utils.ErrNotImplemented
}

func (i *messageInteractor) FindByID(queueID string, msgID string) (*model.Message, error) {
	return nil, utils.ErrNotImplemented
}

func (i *messageInteractor) Add(queueID string, message *model.Message) (*model.Message, error) {
	queue, err := i.queueServie.FindOrCreate(queueID)

	if err != nil {
		return nil, err
	}

	return queue.Add(message)
}

func (i *messageInteractor) Pop(queueID string, qunatity int) ([]*model.Message, error) {
	queue, err := i.queueServie.FindOrCreate(queueID)

	if err != nil {
		return nil, err
	}

	return queue.Pop(1)
}
