package mongodb

import (
	"github.com/totemcaf/quongo/main/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strings"
)

type Database struct {
	mongo  *mgo.Session
	db     *mgo.Database
	dbName string
	log    *log.Logger
}

func NewDatabase(urls string, user string, pwd string, dbName string, log *log.Logger) (*Database, error) {
	dialInfo, err := mgo.ParseURL(urls)

	if err != nil {
		return nil, err
	}

	if user != "" {
		dialInfo.Username = user
		dialInfo.Password = pwd
	}

	log.Printf("Using mongodb as user '%v' on %v", dialInfo.Username, strings.Join(dialInfo.Addrs, ","))
	session, err := mgo.DialWithInfo(dialInfo)

	if err != nil {
		return nil, err
	}

	mgo.SetLogger(log)

	mgo.SetDebug(true)

	db := session.DB(dbName)

	return &Database{session, db, dbName, log}, nil
}

func (d *Database) Close() {
	d.mongo.Close()
}

func (d *Database) Table(tableName string) *mgo.Collection {
	return d.mongo.DB(d.dbName).C(tableName) // TODO put DB() into structure
}

// Get the underlying table for a given queue
func (d *Database) QueueTable(queueName string) *mgo.Collection {
	return d.mongo.DB(d.dbName).C("q-" + queueName)
}

func (d *Database) CreateQueueTable(queueName string) error {
	queue := d.Table(queueName)

	// im.ensure(Index(Seq("ack" -> IndexType.Hashed), unique = true, sparse = true))
	// im.ensure(Index(Seq("visibility" -> IndexType.Ascending), unique = false, sparse = false))

	indexes := []mgo.Index{mgo.Index{
		Key: []string{"hashed:ack"},

		Unique:     true,
		DropDups:   true,
		Background: false, // See notes.
		Sparse:     true,
	},

		mgo.Index{
			Key:        []string{"visibility"},
			Unique:     false,
			DropDups:   false,
			Background: false, // See notes.
			Sparse:     false,
		},

		mgo.Index{
			Key:        []string{"programmed"},
			Unique:     false,
			DropDups:   false,
			Background: false, // See notes.
			Sparse:     false,
		},
	}

	for _, index := range indexes {

		if err := queue.EnsureIndex(index); err != nil {
			return err
		}
	}

	return nil
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
	pipe := r.QueueTable(queueName).Pipe([]bson.M{
		{"$project": bson.M{"v": bson.M{"$lt": []string{"$visible", utils.NowStr()}}}},
		{"$group": bson.M{"_id": "$v", "count": bson.M{"$sum": 1}}},
	})

	counters := make([]VisibleCount, 2)

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
