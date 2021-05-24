package core

import (
	"github.com/go-yaml/yaml"
	"io"
	"os"
)

type Config struct {
	LogLevel       string        `yaml:"log_level"`
	RequestTimeout int           `yaml:"requests_timeout_in_seconds"`
	SentinelList   []string      `yaml:"sentinel_list"`
	DbList         map[string]Db `yaml:"db_list"`
	GraylogHost    string        `yaml:"graylog.host"`
	GraylogPort    string        `yaml:"graylog.port"`
}

type Db struct {
	DbName    string `yaml:"dbname"`
	LocalPort int    `yaml:"local_port"`
}

var config *Config

func NewConfig(configPath string) (*Config, error) {
	cfg := &Config{
		RequestTimeout: 5,
	}

	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer closeFile(f)

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		return nil, err
	}

	config = cfg

	return config, nil
}

func closeFile(f io.Closer) {
	_ = f.Close()
}

func GetConfig() *Config {
	return config
}
