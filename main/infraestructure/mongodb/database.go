package mongodb

import (
  "context"
  "github.com/mongodb/mongo-go-driver/bson"
  "github.com/mongodb/mongo-go-driver/core/connstring"
  "github.com/mongodb/mongo-go-driver/mongo"
  "github.com/totemcaf/quongo/main/utils"
  "log"
  "strings"
)

type Database struct {
	mongo  *mongo.Client
	db     *mongo.Database
	dbName string
	log    *log.Logger
}

func NewDatabase(urls string, user string, pwd string, dbName string, log *log.Logger) (*Database, error) {
  cs, err := connstring.Parse(urls)

	if err != nil {
		return nil, err
	}

	if user != "" {
		cs.Username = user
		cs.Password = pwd
	}

	log.Printf("Using mongodb as user '%v' on %v", cs.Username, strings.Join(cs.Hosts, ","))

	client, err := mongo.NewClientFromConnString(cs)

	if err != nil {
		return nil, err
	}

	// TODO Debug, Logger

	db := client.Database(dbName)   // TODO Set Options !!!

	return &Database{client, db, dbName, log}, nil
}

func (d *Database) Close() {
	d.mongo.Disconnect(context.Background())
}

func (d *Database) Table(tableName string) *mongo.Collection {
	return d.db.Collection(tableName)   // TODO Set Options !!!
}

// Get the underlying table for a given queue
func (d *Database) QueueTable(queueName string) *mongo.Collection {
	return d.db.Collection("q-" + queueName)  // TODO Set Options !!!
}

func (d *Database) CreateQueueTable(queueName string) error {
	queue := d.Table(queueName)

	indexView := queue.Indexes()

	indexNames, err := indexView.CreateMany(context.Background(), []mongo.IndexModel{
      {
        Keys: bson.NewDocument(bson.EC.Int32("ack", 1)),  // TODO It was 'hashed'
        Options: mongo.NewIndexOptionsBuilder().
          Unique(true).Background(false).Sparse(true).Build(),  // TODO DropDups
      },
      {
        Keys: bson.NewDocument(bson.EC.String("visibility", "1")),
        Options: mongo.NewIndexOptionsBuilder().
          Unique(false).Background(false).Sparse(false).Build(),  // TODO DropDups
      },
      {
        Keys: bson.NewDocument(bson.EC.String("programmed", "1")),
        Options: mongo.NewIndexOptionsBuilder().
          Unique(false).Background(false).Sparse(false).Build(),  // TODO DropDups
      },
    },
	)

	if err == nil {
	  log.Printf("Indexes created for queue '%v': %v", queueName, strings.Join(indexNames, ", "))
  }

	// im.ensure(Index(Seq("ack" -> IndexType.Hashed), unique = true, sparse = true))
	// im.ensure(Index(Seq("visibility" -> IndexType.Ascending), unique = false, sparse = false))
	return err
}

type VisibleCount struct {
	Visible bool `bson:"_id"`
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
func (r *Database) Stats(queueName string) (hidden int, visible int) {
  pipeline := bson.NewArray(bson.VC.Document(
    bson.NewDocument(
      bson.EC.SubDocumentFromElements("$project",
        bson.EC.SubDocumentFromElements("v",
          bson.EC.Array("$lt",
            bson.NewArray(bson.VC.Document(
              bson.NewDocument(bson.EC.String("$visible", utils.NowStr())),
            )),
          ),
        ),
      ),
      bson.EC.SubDocumentFromElements("$group",
        bson.EC.String("_id", "$v"),
        bson.EC.SubDocumentFromElements("count", bson.EC.Int32("$sum", 1)),
      ),
    ),
  ))

  pipe, err := r.QueueTable(queueName).Aggregate(context.Background(), pipeline)

  if err != nil {
    hidden = -1
    visible = -1
  	return
	}

	counters := make([]VisibleCount, 2)

	pipe.Decode(&counters)

	for _, c := range counters {
		if c.Visible {
			visible = c.Count
		} else {
			hidden = c.Count
		}
	}
  return
}
