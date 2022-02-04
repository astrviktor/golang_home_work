package config

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type CalendarConfig struct {
	Logger     LoggerConf
	Storage    StorageConf
	HTTPServer HTTPServerConf
	GRPCServer GRPCServerConf
}

type LoggerConf struct {
	Level      int    `yaml:"level"`
	TimeFormat string `yaml:"timeformat"`
}

type StorageConf struct {
	Mode string `yaml:"mode"`
	DSN  string `yaml:"dsn"`
}

type HTTPServerConf struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type GRPCServerConf struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func NewCalendarConfig(configFile string) CalendarConfig {
	var config CalendarConfig

	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		log.Println(err)
		log.Println("using default calendar config")
		return DefaultConfig()
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Println(err)
		log.Println("using default calendar config")
		return DefaultConfig()
	}

	return config
}

func DefaultConfig() CalendarConfig {
	return CalendarConfig{
		LoggerConf{Level: 1, TimeFormat: "2006-01-02T15:04:05Z07:00"},
		StorageConf{Mode: "in-memory", DSN: "postgres://user:password123@localhost:5432/calendar"},
		HTTPServerConf{Host: "", Port: "8888"},
		GRPCServerConf{Host: "", Port: "9999"},
	}
}
