package config

import (
	"encoding/json"
	"io/ioutil"

	"go.uber.org/zap"
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
	Logging zap.Config
}

//TODO add config override with environment variables

//Init the config with json file path
func Init(filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	Instance = new(Config)
	if err := json.Unmarshal(data, Instance); err != nil {
		panic(err)
	}
}
