package clients

import (
	"sync"
)

type Params struct {
	AmpqUri               string
	MongoUri              string
	OmdbUri               string
	TwitterConsumerKey    string
	TwitterConsumerSecret string
	TwitterAccessToken    string
	TwitterAccessSecret   string
}

var (
	initOnce sync.Once
	ampq     Ampq
	mongoDb  Mongo
	omdb     Omdb
	twitter  Twitter
)

func Init(p *Params) {
	initOnce.Do(func() {
		ampq.uri = p.AmpqUri
		mongoDb.uri = p.MongoUri
		omdb.uri = p.OmdbUri
		twitter.accessSecret = p.TwitterAccessSecret
		twitter.accessToken = p.TwitterAccessToken
		twitter.consumerKey = p.TwitterConsumerKey
		twitter.consumerSecret = p.TwitterConsumerSecret
	})
}

func GetOmdb() Omdb {
	return omdb
}

func GetMongoDb() *Mongo {
	return mongoDb.Get()
}

func GetTwitter() *Twitter {
	return twitter.Get()
}

func GetRabbitMq() *Ampq {
	return ampq.Get()
}

//import (
//	"encoding/json"
//	"log"
//	"os"
//)
//
//type ClientsConfig struct {
//	AmpqUri               string
//	MongoUri              string
//	OmdbUri               string
//	TwitterConsumerKey    string
//	TwitterConsumerSecret string
//	TwitterAccessToken    string
//	TwitterAccessSecret   string
//}
//
//var config ClientsConfig
//
//func init() {
//	config = ClientsConfig{}
//}
//
//func Init(configPath string) {
//	file, err := os.Open(configPath)
//	defer file.Close()
//	if err != nil {
//		log.Println("Cannot run grabber without config file: ", err)
//		return
//	}
//
//	decoder := json.NewDecoder(file)
//	err = decoder.Decode(&config)
//	if err != nil {
//		log.Println("Config file should be a valid json: ", err)
//		return
//	}
//}
//
func failOnError(err error, msg string) {

}
