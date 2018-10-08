package grabber

import "log"

type MainConfig struct {
	ApplicationName       string
	ApplicationEnv        string
	ApplicationConfigPath string
	ApplicationTasks      []string
}

var config MainConfig

func Init(config MainConfig) {
	config = config
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
