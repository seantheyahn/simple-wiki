package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

//Instance config instance
var Instance *Config

//Config --
type Config struct {
	Server   serverConfig
	Db       dbConfig
	Redis    redisConfig
	RootUser rootUserConfig
}

type dbConfig struct {
	ConnectionURI string
	MaxOpenConns  int
	MaxIdleConns  int
}

type redisConfig struct {
	ConnectionURI string
}

type serverConfig struct {
	ListenAddress string
}

type rootUserConfig struct {
	Username string
	Password string
}

//Init the config with json file path
func Init(filename string) {
	log.Println("initialzing config")
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	Instance = new(Config)
	if err := json.Unmarshal(data, Instance); err != nil {
		panic(err)
	}
}
