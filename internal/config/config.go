package config

import (
	"errors"
	"os"

	"emailservice/pkg/logger"

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
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"smtp"`
	Email struct {
		SecretKey string `yaml:"secretKey"`
	}
}

//NewConfig loads the config file
func NewConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		logger.Fatal(err)
	}
	confPath := os.Getenv("EMAIL_CONFIG")
	if confPath == "" {
		return nil, errors.New("EMAIL_CONFIG env variable is not set")
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
