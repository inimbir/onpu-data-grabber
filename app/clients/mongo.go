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

type Tweet struct {
	Id    bson.ObjectId `bson:"_id"`
	Group string        `bson:"group"`
	Ids   []string      `bson:"ids"`
}

type HashTag struct {
	Id       bson.ObjectId `bson:"_id"`
	Group    string        `bson:"group"`
	HashTags []Tag         `bson:"hashtags"`
}

type Tag struct {
	Name   string `bson:"name";json:"name"`
	Status int    `bson:"status";json:"status"`
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
	//func (m Mongo) ExistsHashTag(group string, name string, statuses []int) (exists bool, err error) {
	coll := m.db.C(HashTagsCollection)
	n := 0
	if n, err = coll.Find(bson.M{
		"group":         group,
		"hashtags.name": name,
	}).Count(); err != nil {
		return false, err
	}
	exists = n > 0
	return
}

func (m Mongo) UpdateHashTagStatus(group string, name string, status int) (err error) {
	coll := m.db.C(HashTagsCollection)
	if err = coll.Update(bson.M{
		"group":         group,
		"hashtags.name": name,
	}, bson.M{"$set": bson.M{"hashtags.$.status": status}}); err != nil {
		return err
	}
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
	change := bson.M{"$addToSet": bson.M{"hashtags": bson.M{"$each": hashTag.HashTags}}}
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
	for _, tag := range product.HashTags {
		hashtags = append(hashtags, tag.Name)
	}
	return
}

func (m Mongo) GetTagsByGroup(group string, tagType int) (hashtags []Tag, err error) {
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
		return hashtags, fmt.Errorf("cannot fetch data to model")
	}
	if len(product.HashTags) == 0 {
		return hashtags, fmt.Errorf("model not exists or tags are empty")
	}
	return product.HashTags, err
}

func (m Mongo) ExistsTweet(group string, id int64) (exists bool, err error) {
	coll := m.db.C(TweetsCollection)
	n := 0
	if n, err = coll.Find(bson.M{
		"group": group,
		"ids":   bson.M{"$in": []int64{id}},
	}).Count(); err != nil {
		return false, err
	}
	exists = n > 0
	return
}

func (m Mongo) InsertTweet(group string, id string) (err error) {
	coll := m.db.C(TweetsCollection)
	tweet := &Tweet{
		Id:    bson.NewObjectId(),
		Group: group,
		Ids:   []string{id}}
	if err = coll.Insert(tweet); err != nil {
		fmt.Println(err)
	}
	return
}

func (m Mongo) AddTweetIdToGroup(group string, ids []int64) (err error) {
	coll := m.db.C(TweetsCollection)
	change := bson.M{"$addToSet": bson.M{"ids": bson.M{"$each": ids}}}
	if _, err := coll.Upsert(bson.M{
		"group": group,
	}, change); err != nil {
		return err
	}
	return
}

//find by name and status
//	products := []Group{}
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
//products := []Group{}
//coll.Find(bson.M{
//	"group": "tt65464",
//	"hashtags.status": bson.M{
//		"$eq": 0,
//	},
//}).All(&products)
//for _, p := range products {
//	fmt.Println(p.Group, p.HashTags)
//}
