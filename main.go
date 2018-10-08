package main

import (
	"encoding/json"
	"flag"
	"github.com/inimbir/onpu-data-grabber/clients"
	"github.com/inimbir/onpu-data-grabber/grabber"
	"log"
	"os"
)

const (
	ApplicationName           = "onpu-data-grabber"
	ApplicationEnv            = "develop"
	ApplicationConfigFilePath = "config/main.json"
	ApplicationTasksFilePath  = "config/tasks.json"
)

type Config struct {
	MainConfig    grabber.MainConfig
	ClientsConfig clients.ClientsConfig
}

func main() {
	err := initConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(clients.GetSeriesId("American Vandal"))

	//	config, err := initConfig();
	//log.Println(config.ClientsConfig.OmdbUri, err)
	//if err != nil {
	//	log.Println("Error loading config: ", err)
	//	return
	//}
	//clients.Init(config.ClientsConfig)
	//grabber.Init(config.MainConfig)
	//grabber.Run();
}

func initConfig() (err error) {
	var tasksFilePath string
	c := &Config{
		MainConfig: grabber.MainConfig{
			ApplicationName: ApplicationName,
		},
	}

	flag.StringVar(&c.MainConfig.ApplicationEnv, "env", ApplicationEnv, "Environment of application")
	flag.StringVar(&c.MainConfig.ApplicationConfigPath, "env-file-path", ApplicationConfigFilePath, "Path to file with application configuration (login, secrets)")
	flag.StringVar(&tasksFilePath, "tasks-file-path", ApplicationTasksFilePath, "Path to file with application configuration (login, secrets)")
	flag.Parse()
	tasksFile, err := os.Open(tasksFilePath)
	defer tasksFile.Close()
	if err != nil {
		log.Println("Cannot run grabber without tasks file: ", err)
		return
	}

	tasksDecoder := json.NewDecoder(tasksFile)
	err = tasksDecoder.Decode(&c.MainConfig)
	if err != nil {
		log.Println("Tasks file should be a valid json: ", err)
		return
	}

	clients.Init(c.MainConfig.ApplicationConfigPath)
	grabber.Init(c.MainConfig)

	return nil
}
