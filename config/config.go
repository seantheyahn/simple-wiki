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
	Server struct {
		ListenAddress string
		CookieSecret  string
		CSRFSecret    string
		IPRateLimiter struct {
			ReqPerSecond int
			ReqBurst     int
		}
	}
	Db struct {
		ConnectionURI string
		MaxOpenConns  int
		MaxIdleConns  int
	}
	Redis struct {
		ConnectionURI string
	}
	RootUser struct {
		Username string
		Password string
	}
}

//TODO add config override with environment variables

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
