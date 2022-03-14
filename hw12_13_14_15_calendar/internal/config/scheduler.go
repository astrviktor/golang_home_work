package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type SchedulerConfig struct {
	AMQPScheduler AMQPSchedulerConf `yaml:"amqpScheduler"`
	Storage       StorageConf       `yaml:"storage"`
}

type AMQPSchedulerConf struct {
	URI          string `yaml:"uri"`
	Exchange     string `yaml:"exchange"`
	RepeatSecond int    `yaml:"repeatSecond"`
}

func NewSchedulerConfig(configFile string) SchedulerConfig {
	var config SchedulerConfig

	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		log.Println(err)
		log.Println("using default scheduler config")
		return DefaultSchedulerConfig()
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Println(err)
		log.Println("using default scheduler config")
		return DefaultSchedulerConfig()
	}

	return config
}

func DefaultSchedulerConfig() SchedulerConfig {
	return SchedulerConfig{
		AMQPSchedulerConf{"amqp://guest:guest@localhost:5672/", "events", 5},
		StorageConf{Mode: "sql", DSN: "postgres://user:password123@localhost:5432/calendar"},
	}
}
