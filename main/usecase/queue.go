package usecase

import (
	"github.com/totemcaf/quongo/main/model"
	"github.com/totemcaf/quongo/main/utils"
)

type queue struct {
}

func NewQueueInteractor() *queue {
	return &queue{}
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

func (i *queue) IsQueueNameValid(queueID string) bool {
	return false
}

func (i *queue) Add(queue model.Queue) (*model.Queue, error) {
	return nil, utils.ErrNotImplemented
}

func (i *queue) Update(queue model.Queue) (*model.Queue, error) {
	return nil, utils.ErrNotImplemented
}
