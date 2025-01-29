package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

var Conf Config

type Config struct {
	Database Database `yaml:"database"`
	Cache    Cache    `yaml:"cache"`
	Auth     Auth     `yaml:"auth"`
}

type Database struct {
	Dsn                   string `yaml:"dsn"`
	MaxIdle               int    `yaml:"max_idle"`
	MaxOpenConnection     int    `yaml:"max_open_connection"`
	MaxLifeTimeConnection int64  `yaml:"max_life_time_connection"`
}

type Cache struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type Auth struct {
	SecretKey  string `yaml:"secret_key"`
	HashSecret string `yaml:"hash_secret"`
}

func init() {
	fl, err := os.ReadFile("./config/config.yml")
	if err != nil {
		log.Fatal(err)
	}
	if err = yaml.Unmarshal([]byte(os.ExpandEnv(string(fl))), &Conf); err != nil {
		log.Fatal(err)
	}
}
