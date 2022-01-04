package main

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Logger  LoggerConf
	Storage StorageConf
	Server  HTTPServerConf
}

type LoggerConf struct {
	Level      string `yaml:"level"`
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

func NewConfig(configFile string) Config {
	var config Config

	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Println(err)
		log.Println("using default config...")
		return DefaultConfig()
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Println(err)
		log.Println("using default config...")
		return DefaultConfig()
	}

	return config
}

func DefaultConfig() Config {
	return Config{
		LoggerConf{Level: "INFO", TimeFormat: "2006-01-02T15:04:05Z07:00"},
		StorageConf{Mode: "in-memory", DSN: "postgres://user:password123@localhost:5432/calendar?sslmode=require"},
		HTTPServerConf{Host: "", Port: "8888"},
	}
}

// TODO
