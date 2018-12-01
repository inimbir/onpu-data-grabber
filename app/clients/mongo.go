package clients

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"sync"
)

type Mongo struct {
	uri     string
	session *mgo.Session
	db      *mgo.Database
}

var (
	instance      *Mongo
	initMongoOnce sync.Once
)

const TweetsCollection = "tweets"
const HashTagsCollection = "hashtags"

func (m Mongo) Get() *Mongo {
	initMongoOnce.Do(func() {
		var err error
		if m.session, err = mgo.Dial(m.uri); err != nil {
			log.Fatalf("Could not create mongo db client")
		}
		//defer m.session.Close()
		m.db = m.session.DB("heroku_p30gwwqt")
		instance = &m
	})
	return instance
}

type HashTag struct {
	Id       bson.ObjectId `bson:"_id"`
	Group    string        `bson:"group"`
	HashTags []Tag         `bson:"hashtags"`
}

type Tag struct {
	Name   string `bson:"name"`
	Status int    `bson:"status"`
}

func (m Mongo) InsertHashTag(group string, name string, status int) (err error) {
	coll := m.db.C(HashTagsCollection)
	hashTag := &HashTag{Id: bson.NewObjectId(), Group: group, HashTags: []Tag{{
		Name:   name,
		Status: status,
	}}}
	if err = coll.Insert(hashTag); err != nil {
		fmt.Println(err)
	}
	return
}

func (m Mongo) ExistsHashTag(group string, name string) (exists bool, err error) {
	coll := m.db.C(HashTagsCollection)
	n := 0
	if n, err = coll.Find(bson.M{
		"group":         group,
		"hashtags.name": name,
	}).Count(); err != nil {
		fmt.Println(err)
	}
	exists = n > 0
	return
}

func (m Mongo) BulkInsert(group string, tags []string, status int) (err error) {
	coll := m.db.C(HashTagsCollection)
	hashTag := &HashTag{}
	for _, tag := range tags {
		hashTag.HashTags = append(hashTag.HashTags, Tag{
			Name:   tag,
			Status: status,
		})
	}
	change := bson.M{"$push": bson.M{"hashtags": bson.M{"$each": hashTag.HashTags}}}
	if err = coll.Update(bson.M{
		"group": group,
	}, change); err != nil {
		fmt.Println(err)
	}
	return
}

func (m Mongo) GetHashTagsByGroup(group string, tagType int) (hashtags []string, err error) {
	coll := m.db.C(HashTagsCollection)
	product := &HashTag{}
	var condition bson.M
	if tagType == -1 {
		condition = bson.M{
			"$match": bson.M{
				"group": group,
			},
		}
	} else {
		condition = bson.M{
			"$match": bson.M{
				"group":           group,
				"hashtags.status": tagType,
			},
		}
	}

	pipeline := []bson.M{
		{"$unwind": "$hashtags"},
		condition,
		{
			"$group": bson.M{
				"_id":      "$_id",
				"hashtags": bson.M{"$push": "$hashtags"},
			},
		},
	}
	pipe := coll.Pipe(pipeline)
	if pipe.One(product); err != nil {
		log.Println("no series")
		return hashtags, fmt.Errorf("no series")
	}
	log.Println(product)

	for _, tag := range product.HashTags {
		hashtags = append(hashtags, tag.Name)
	}
	log.Println(hashtags)

	return
}

//find by name and status
//	products := []HashTag{}
//	coll.Find(bson.D{
//		{"hashtags", bson.D{
//			{"name", "jddournal1"},
//			{"status", 0},
//		}},
//	}).All(&products)
//	for _, p := range products{
//		fmt.Println(p.Group, p.HashTags)
//	}

//find all by status
//products := []HashTag{}
//coll.Find(bson.M{
//	"group": "tt65464",
//	"hashtags.status": bson.M{
//		"$eq": 0,
//	},
//}).All(&products)
//for _, p := range products {
//	fmt.Println(p.Group, p.HashTags)
//}
