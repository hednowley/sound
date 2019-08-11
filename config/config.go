package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Port                     string
	ArtPath                  string `yaml:"art path"`
	Db                       string
	LogConfig                string   `yaml:"log config"`
	IgnoredArticles          []string `yaml:"ignored articles"`
	Users                    []User
	Secret                   string
	AccessControlAllowOrigin string               `yaml:"access control allow origin"`
	WebsocketTicketExpiry    int                  `yaml:"websocket ticket expiry"`
	BeetsProviders           []BeetsProvider      `yaml:"beets"`
	FileSystemProviders      []FileSystemProvider `yaml:"filesystem"`
}

type BeetsProvider struct {
	Database string
	Name     string
}

type FileSystemProvider struct {
	Path       string
	Name       string
	Extensions []string
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
