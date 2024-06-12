package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

// Config Configuration struct to hold application settings
type Config struct {
	MySQL MySQLConfig
	Redis RedisConfig
}

// MySQLConfig holds MySQL database connection settings
type MySQLConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

// RedisConfig holds Redis connection settings
type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

func GetConfig(filepath string) (conf *Config) {
	yamlFile, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal("yamlFile.Get err   #%v " + err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		log.Fatal("Unmarshal: %v " + err.Error())
	}
	return conf
}
