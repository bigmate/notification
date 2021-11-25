package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

//Config is the app config struct
type Config struct {
	AppName string `yaml:"app_name"`
	Logger  struct {
		Level string `yaml:"level"`
	}
	Smtp struct {
		Host      string `yaml:"host"`
		Port      int    `yaml:"port"`
		Username  string `yaml:"username"`
		Password  string `yaml:"password"`
		SecretKey string `yaml:"secretKey"`
		Sender    string `yaml:"sender"`
	} `yaml:"smtp"`
	RabbitMQ struct {
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DSN      string `yaml:"dsn"`
	}
}

//NewConfig loads the config file
func NewConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	confPath := os.Getenv("NOTIFICATION_CONFIG")
	if confPath == "" {
		return nil, errors.New("config: NOTIFICATION_CONFIG env variable is not set")
	}

	conf := &Config{}
	file, err := os.Open(confPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if err = yaml.NewDecoder(file).Decode(conf); err != nil {
		return nil, err
	}

	return conf, nil
}
