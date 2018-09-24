package mgo

import (
	"github.com/totemcaf/quongo/main/model"
	"github.com/totemcaf/quongo/main/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type MessageRepositoryProvider struct {
	db *Database
}

func NewMongodbRepositoryProvider(db *Database) *MessageRepositoryProvider {
	return &MessageRepositoryProvider{db}
}

func (p *MessageRepositoryProvider) ForQueue(queueId string) model.MessageRepository {
	return NewMongodbRepository(queueId, p.db)
}

type MessageRepository struct {
	queueId string
	db      *Database
	table   *mgo.Collection
}

func NewMongodbRepository(queueId string, db *Database) *MessageRepository {
	return &MessageRepository{queueId, db, db.QueueTable(queueId)}
}

// Return a selector of visible messages
func isVisible() bson.M {
	return bson.M{"visible": bson.M{"$lte": time.Now()}}
}

// A sort expression by the original visible time
const byProgrammed = "programmed"

/**
 * Find at most 'limit' visible messages from 'offset' sorted by the original visible time
 */
func (r *MessageRepository) Find(offset int, limit int) ([]model.Message, error) {
	var result []model.Message

	err := r.table.Find(isVisible()).Sort(byProgrammed).Skip(offset).Limit(limit).All(&result)

	return result, err
}

/**
 * Return the message with the provided id or nil.
 */
func (r *MessageRepository) FindById(mid string) (*model.Message, error) {
	var result model.Message
	e1 := r.table.FindId(mid).One(&result)

	return &result, e1
}

/**
 * Returns the first message that is visible and set it as taken (set its ACK value, and increment the retries)
 * If no message is pending, returns nil
 */
func (r *MessageRepository) Pop(delay time.Duration) (*model.Message, error) {
	change := mgo.Change{
		Update: bson.M{
			"ack":     bson.NewObjectId(),
			"$inv":    bson.M{"retries": 1},
			"visible": utils.FromNow(delay),
		},
		ReturnNew: true,
	}

	var result model.Message
	info, e1 := r.table.Find(isVisible()).Sort(byProgrammed).Apply(change, &result)

	if info.Matched == 0 {
		return nil, e1
	}

	return &result, e1
}

/**
 * Add a message into the queue
 */
func (r *MessageRepository) Add(message *model.Message) (*model.Message, error) {
	return message, r.table.Insert(message)
}
