package app

import "github.com/kelseyhightower/envconfig"

type Config struct {
	ServerPort       string `default:":8080"`
	DatabaseHost     string `default:"localhost"`
	DatabasePort     string `default:"54321"`
	DatabaseUsername string `default:"todo"`
	DatabasePassword string `default:"todo"`
	DatabaseSchema   string `default:"todo"`
	Url              string `default:"http://localhost:8080"`
}

func CreateConfig() (*Config, error) {
	c := &Config{}

	err := envconfig.Process("todo", c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
