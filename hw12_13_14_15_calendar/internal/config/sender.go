package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type SenderConfig struct {
	AMQPSender AMQPSenderConf
}

type AMQPSenderConf struct {
	URI   string `yaml:"uri"`
	Queue string `yaml:"queue"`
	Retry int    `yaml:"retry"`
}

func NewSenderConfig(configFile string) SenderConfig {
	var config SenderConfig

	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		log.Println(err)
		log.Println("using default sender config")
		return DefaultSenderConfig()
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Println(err)
		log.Println("using default sender config")
		return DefaultSenderConfig()
	}

	return config
}

func DefaultSenderConfig() SenderConfig {
	return SenderConfig{
		AMQPSenderConf{"amqp://guest:guest@localhost:5672/", "events", 5},
	}
}
