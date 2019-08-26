package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Port                     int
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

// NewConfig creates config from the YAML at the provided path.
// It's defined as a nullary function so we can inject the path before passing
// the constructor to FX.
func NewConfig(path string) func() (*Config, error) {
	return func() (*Config, error) {
		yamlData, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		var c Config
		err = yaml.Unmarshal(yamlData, &c)
		return &c, err
	}
}
