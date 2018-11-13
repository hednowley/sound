package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Port            string
	Path            string
	ArtPath         string `yaml:"art path"`
	Db              string
	LogConfig       string `yaml:"log config"`
	Extensions      []string
	IgnoredArticles []string `yaml:"ignored articles"`
	Users           []User
	BeetsDB         string `yaml:"beets db"`
	Secret          string
}

type User struct {
	Username string
	Password string
	Email    string
}

func NewConfig() (*Config, error) {
	yamlData, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}

	var c Config
	err = yaml.Unmarshal(yamlData, &c)
	return &c, err
}
