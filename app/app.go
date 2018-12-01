package app

import (
	"encoding/json"
	"github.com/inimbir/onpu-data-grabber/app/clients"
	"github.com/inimbir/onpu-data-grabber/app/grabber"
	"log"
	"os"
)

type Context struct {
	Tasks []string
}

const (
	ApplicationName           = "onpu-data-grabber"
	ApplicationConfigFilePath = "config/main.json"
)

func Run() {
	Init()
	grabber.Init()
}

func Init() {
	var (
		file   = &os.File{}
		config = &clients.Params{}
		err    error
	)
	if file, err = os.Open(ApplicationConfigFilePath); err != nil {
		log.Fatalf("Cannot run grabber without config file: %s", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(config); err != nil {
		log.Fatalf("Config file should be a valid json: %s", err)
		return
	}
	clients.Init(config)
}

//type MainConfig struct {
//	ApplicationName       string
//	ApplicationEnv        string
//	ApplicationConfigPath string
//	TasksFilePath string
//	ApplicationTasks      []string
//}
//
//type Config struct {
//	MainConfig    MainConfig
//ClientsConfig clients.ClientsConfig
//}

//var config *Config

//var (
//	tasksFile *os.File
//	err error
//)
//config = &Config{
//	MainConfig: MainConfig{
//		ApplicationName: ApplicationName,
//	},
//}
//flag.StringVar(&c.MainConfig.ApplicationEnv, "env", ApplicationEnv, "Environment of application")
//flag.StringVar(&config.MainConfig.ApplicationConfigPath, "env-file-path", ApplicationConfigFilePath, "Path to file with application configuration (login, secrets)")
//flag.StringVar(&config.MainConfig.TasksFilePath, "tasks-file-path", ApplicationTasksFilePath, "Path to file with application configuration (login, secrets)")
//flag.Parse()
//if tasksFile, err = os.Open(config.MainConfig.TasksFilePath); err != nil {
//	log.Println("Cannot run grabber without tasks file: ", err)
//	return
//}
//defer tasksFile.Close()
//
//
//tasksDecoder := json.NewDecoder(tasksFile)
//err = tasksDecoder.Decode(&c.MainConfig)
//if err != nil {
//	log.Println("Tasks file should be a valid json: ", err)
//	return
//}
//
//clients.Init(c.MainConfig.ApplicationConfigPath)
//grabber.Init(c.MainConfig)
//
//return nil
