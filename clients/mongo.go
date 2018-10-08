package clients

import (
	"context"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"log"
	"testing"
)

var db mongo.Database

func __init() {
	client, err := mongo.Connect(context.Background(), config.MongoUri, nil)
	failOnError(err, "Could not connect to db")

	db = *client.Database("documentation_examples")
}
func GetTags(collection string) {
	//tag, err := db.Collection(collection).Find(
	//	context.Background(),
	//	bson.NewDocument(
	//		bson.EC.String("status", "A"),
	//		bson.EC.SubDocumentFromElements("qty",
	//			bson.EC.Int32("$lt", 30),
	//		),
	//	))
	//failOnError(err, "failed to execute search")
	//return tag.ID()
}

func IsTweetProcessed(collection string, id string) {
	//tag, err := db.Collection(collection).Find(
	//	context.Background(),
	//	bson.NewDocument(
	//		bson.EC.String("status", "A"),
	//		bson.EC.SubDocumentFromElements("qty",
	//			bson.EC.Int32("$lt", 30),
	//		),
	//	))
}

func Insert(collection string, id string) {

}

func InitCollection(collection string, id string) {

	coll := db.Collection(collection)
	result, err := coll.InsertOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("item", "canvas"),
			bson.EC.Int32("qty", 100),
			bson.EC.ArrayFromElements("tags",
				bson.VC.String("cotton"),
			),
			bson.EC.SubDocumentFromElements("size",
				bson.EC.Int32("h", 28),
				bson.EC.Double("w", 35.5),
				bson.EC.String("uom", "cm"),
			),
		))
	log.Println(result)
	log.Println(err)

}

func InsertExamples(t *testing.T, db *mongo.Database, collection string, id string) {
	_, err := db.RunCommand(
		context.Background(),
		bson.NewDocument(bson.EC.Int32("dropDatabase", 1)),
	)

	failOnError(err, "failed to create collection")

	coll := db.Collection(collection)

	{
		// Start Example 1

		result, err := coll.InsertOne(
			context.Background(),
			bson.NewDocument(
				bson.EC.String("item", "canvas"),
				bson.EC.Int32("qty", 100),
				bson.EC.ArrayFromElements("tags",
					bson.VC.String("cotton"),
				),
				bson.EC.SubDocumentFromElements("size",
					bson.EC.Int32("h", 28),
					bson.EC.Double("w", 35.5),
					bson.EC.String("uom", "cm"),
				),
			))

		// End Example 1

		log.Println(err)
		log.Println(result)
		//require.NoError(t, err)
		//require.NotNil(t, result.InsertedID)
	}

	{
		// Start Example 2

		cursor, err := coll.Find(
			context.Background(),
			bson.NewDocument(bson.EC.String("item", "canvas")),
		)

		// End Example 2

		//require.NoError(t, err)
		//requireCursorLength(t, cursor, 1)

		log.Println(err)
		log.Println(cursor)
	}

}
