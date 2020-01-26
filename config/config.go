package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"go.uber.org/zap"
)

//Instance config instance
var Instance *Config

//Config --
type Config struct {
	Server struct {
		ListenAddress   string
		CookieSecret    string
		CSRFSecret      string
		DevelopmentMode bool
		IPRateLimiter   struct {
			ReqPerSecond int
			ReqBurst     int
		}
	}
	DB struct {
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
	Logging struct {
		Level zap.AtomicLevel
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func extractMap(key string, val string, parent map[string]interface{}) {
	if !strings.Contains(key, "__") {
		var x interface{}
		checkError(json.Unmarshal([]byte(val), &x))
		parent[key] = x
		return
	}

	sp := strings.SplitN(key, "__", 2)
	k := sp[0]
	var m map[string]interface{}
	x, ok := parent[k]
	if ok {
		m, ok = x.(map[string]interface{})
		if !ok {
			m = make(map[string]interface{})
			parent[k] = m
		}
	} else {
		m = make(map[string]interface{})
		parent[k] = m
	}

	extractMap(sp[1], val, m)
}

func getEnvMap(prefix string) map[string]interface{} {
	m := make(map[string]interface{})
	envs := os.Environ()
	for _, env := range envs {
		if !strings.HasPrefix(env, prefix) {
			continue
		}
		sp := strings.SplitN(strings.TrimPrefix(env, prefix), "=", 2)
		key := sp[0]
		val := sp[1]
		extractMap(key, val, m)
	}
	return m
}

//Init intis the config with json file path and override with env:
// override usage: CONF__MyKey__MyInnerKey="MyJsonValue"
// example 1: export CONF__Server__ListenAddress=\":8080\"
// example 2: export CONF__DB__MaxOpenConns=10
func Init(filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	Instance = new(Config)

	//load from file
	checkError(json.Unmarshal(data, Instance))

	j, err := json.Marshal(getEnvMap("CONF__"))
	checkError(err)

	//override with env
	checkError(json.Unmarshal(j, Instance))
}
