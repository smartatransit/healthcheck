package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Options struct {
		PollTimeSeconds int `yaml:"poll_time_seconds"`
		TimeoutSeconds  int `yaml:"timeout_seconds"`
	} `yaml:"options"`
	Services struct {
		Name     string `yaml:"name"`
		Endpoint string `yaml:"endpoint"`
		Enabled  bool   `yaml:"enabled"`
	} `yaml:"services"`
}

func NewConfig(filepath string) (Config, error) {
	config := Config{}
	file, err := os.Open(filepath)
	if err != nil {
		err = errors.Wrapf(err, "failed opening config file %s for reading", filepath)
		return config, err
	}
	err = yaml.NewDecoder(file).Decode(&config)
	err = errors.Wrapf(err, "failed parsing config file %s", file.Name())
	return config, err

}
