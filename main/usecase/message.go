package usecase

import (
	"github.com/totemcaf/quongo/main/model"
	"github.com/totemcaf/quongo/main/utils"
)

type messageInteractor struct {
}

func NewMessageInteractor() *messageInteractor {
	return &messageInteractor{}
}

func (i *messageInteractor) FindAll(queueID string, offset int, limit int) ([]model.Message, error) {
	return nil, utils.ErrNotImplemented
}

func (i *messageInteractor) FindByID(queueID string, msgID string) (*model.Message, error) {
	return nil, utils.ErrNotImplemented
}

func (i *messageInteractor) Add(queueID string, message *model.Message) (*model.Message, error) {
	return nil, utils.ErrNotImplemented
}

func (i *messageInteractor) Pop(queueID string, qunatity int) ([]*model.Message, error) {
	return nil, utils.ErrNotImplemented
}
