package clients

import (
	"encoding/json"
	"log"
	"os"
)

type ClientsConfig struct {
	AmpqUri               string
	MongoUri              string
	OmdbUri               string
	TwitterConsumerKey    string
	TwitterConsumerSecret string
	TwitterAccessToken    string
	TwitterAccessSecret   string
}

var config ClientsConfig

func init() {
	config = ClientsConfig{}
}

func Init(configPath string) {
	file, err := os.Open(configPath)
	defer file.Close()
	if err != nil {
		log.Println("Cannot run grabber without config file: ", err)
		return
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Println("Config file should be a valid json: ", err)
		return
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
