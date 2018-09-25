package mongodb

import (
	"context"
	"errors"
	"github.com/mongodb/mongo-go-driver/bson"
  "github.com/mongodb/mongo-go-driver/mongo/findopt"
	"github.com/mongodb/mongo-go-driver/mongo"

	"github.com/totemcaf/quongo/main/model"
	"log"
)

type Repository struct {
	db        *Database
	queueColl *mongo.Collection
}

func NewQueueRepository(db *Database) *Repository {
	queueColl := db.Table("model.Queue")

	return &Repository{
		db,
		queueColl,
	} // TODO Create indexes
}

/**
 * Report at most 'size' queues starting from 'skip'
 */
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

	// result := make([]*model.Queue, size) // TODO Should be "*" or struct directly?

	cursor, err := r.queueColl.Find(context.Background(), bson.NewDocument(),
	  findopt.Skip(int64(skip)), findopt.Limit(int64(size)),
  )

	var queues []*model.Queue

	for cursor.Next(context.Background()) {
	  queue := model.Queue{}

	  err := cursor.Decode(queue)

	  if err != nil {
	    return nil, err
    }

    queues = append(queues, &queue)
  }

	return queues, err
}

/**
 * Add the queue statistics to the provided queue
 */
func (r *Repository) Complete(queue *model.Queue) *model.QueueWithStats {
	hidden, visible := r.db.Stats(queue.Name)

	return model.NewQueueWithStat(queue, visible+hidden, hidden, visible)
}

/**
 * Return a queue with the provided id (name)
 */
func (r *Repository) FindById(id string) (*model.Queue, error) {
	log.Printf("Finding queue by id: %v", id)
	result := model.Queue{}

	err := r.queueColl.FindId(id).One(&result)

	if err == nil {
		log.Printf("Found: %v", result.Name)
		return &result, nil
	} else {
		log.Printf("Not found")
		return nil, err
	}
}

func (r *Repository) Add(queue *model.Queue) (*model.Queue, error) {
	log.Printf("Inserting queue %v", queue.Name)
	result, err := r.queueColl.InsertOne(context.Background(), queue)

	if err != nil {
		log.Printf("Error:      %v", err)
	}
	return queue, err
}

func (r *Repository) Update(queue *model.Queue) (*model.Queue, error) {
	return queue, nil
}
