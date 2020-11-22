package server

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Port  int         `yaml:"port"`
	Mongo MongoConfig `yaml:"mongo"`
}

type MongoConfig struct {
	Url        string `yaml:"url"`
	Database   string `yaml:"database"`
	Collection string `yaml:"collection"`
}

func ReadConfig(configPath string) (cfg Config, err error) {
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(yamlFile, &cfg)
	return
}
