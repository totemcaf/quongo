package usecase

import (
	"github.com/totemcaf/quongo/main/model"
	"time"
)

type MessageRepositoryProvider interface {
	ForQueue(queueId string) model.MessageRepository
}

type MessageInteractor struct {
	repositoryProvider MessageRepositoryProvider
}

func NewMessageInteractor(repositoryProvider MessageRepositoryProvider) *MessageInteractor {
	return &MessageInteractor{repositoryProvider: repositoryProvider}
}

func (i *MessageInteractor) FindAll(queueId string, offset int, limit int) ([]model.Message, error) {
	return i.repositoryProvider.ForQueue(queueId).Find(offset, limit)
}

func (i *MessageInteractor) FindById(queueId string, msgId string) (*model.Message, error) {
	return i.repositoryProvider.ForQueue(queueId).FindById(msgId)
}

func (i *MessageInteractor) Add(queueId string, message *model.Message) (*model.Message, error) {
	return i.repositoryProvider.ForQueue(queueId).Add(message)
}

func (i *MessageInteractor) Pop(queueId string, delay time.Duration) (*model.Message, error) {
	return i.repositoryProvider.ForQueue(queueId).Pop(delay)
}
