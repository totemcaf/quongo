package queue

import (
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "errors"
  "github.com/totemcaf/quongo/main/message"
)

type Repository struct {
  mongo     *mgo.Session
  db        *mgo.Database
  queueColl *mgo.Collection
  msgRep    *message.Repository
}

func NewQueueRepository(mongo *mgo.Session, dbName string, msgRep *message.Repository) *Repository {
  db := mongo.DB(dbName)
  queueColl := db.C("Queue")

  return &Repository{mongo, db, queueColl, msgRep} // TODO Create indexes
}

func NewQueueStats(name string) *QueueWithStats {
  pQ, _ := NewQueue(name)

  return &QueueWithStats{
    Queue: pQ,
    Stats: QueueStats{10, 2, 3},
  }
}

func (r *Repository) findAll(skip int, size int) (*[] Queue, error) { // TODO Should be *Queue ??
  if skip < 0 {
    return nil, errors.New("Invalid skip, it should be >= 0")
  }

  if size < 0 {
    return nil, errors.New("Invalid size, it shoud be >= 0")
  }

  if size == 0 {
    return &[] Queue{}, nil
  }

  result := make([] Queue, size)

  err := r.queueColl.Find(bson.M{}).Skip(skip).Limit(size).All(&result)

  return &result, err
}

func (this *Repository) complete(queue *Queue) *QueueWithStats {
  hidden, visible := this.msgRep.Stats(queue.Name)

  return NewQueueWithStat(queue, visible + hidden, hidden, visible)
}

func (r *Repository) findById(id string) (*QueueWithStats, error) {
  result := Queue{}

  err := r.queueColl.Find(bson.M{"_id": id}).One(&result)

  if err == nil {
    return r.complete(&result), nil
  } else {
    return nil, err
  }
}

func (r *Repository) add(pQueue *Queue) (*Queue, error) {
  return pQueue, nil
}

func (r *Repository) Update(pQueue *Queue) (*Queue, error) {
  return pQueue, nil
}
