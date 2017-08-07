package message
import (
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "time"
)

type Repository struct {
  mongo   *mgo.Session
  db      *mgo.Database
  dbName  string
}

func NewMsgRepository(mongo *mgo.Session, dbName string) *Repository {
  db := mongo.DB(dbName)
  return &Repository{mongo, db, dbName}
}


func (r *Repository) table(queueName string) *mgo.Collection {
  return r.mongo.DB(r.dbName).C(queueName)
}

type VisibleCount struct {
  Visible bool   `bson:"_id"`
  Count   int
}

/*
db["q-event-d365"].aggregate([
  { $project: {v: { $lt: [ "$visible", ISODate("2017-07-14T15:48:16.182Z")] }}},
  { $group: { _id: "$v", count: {$sum: 1}} }
])

TODO visible and inProcess son distintos !!!
  visible : visible  < ahora
  inProcess: ! visible && ack != null
 */
func (this *Repository) Stats(queueName string) (hidden int, visible int) {
  now := time.Now()

  pipe := this.table("q-" + queueName).Pipe([]bson.M{
    {"$project": bson.M{"v": bson.M{"$lt": [] string {"$visible", now.Format("2006-01-02T15:04:05.000Z")}}}},
    {"$group": bson.M{"_id": "$v", "count": bson.M{"$sum": 1}}},
  })

  counters := make([] VisibleCount, 2)

  pipe.All(&counters)

  for _, c := range counters {
    if c.Visible {
      visible = c.Count
    } else {
      hidden = c.Count
    }
  }

  return
}