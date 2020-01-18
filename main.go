package main

import (
	"flag"
	"fmt"
	"log"
	"sean/wiki/config"
	"sean/wiki/services"
	"strings"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	"github.com/json-iterator/go/extra"
	cmap "github.com/orcaman/concurrent-map"

	_ "github.com/jackc/pgx/stdlib"
)

var userMap cmap.ConcurrentMap = cmap.New()
var idCounter uint64 = 0

type User struct {
	Id   string // `json:"id"`
	Name string // `json:"name"`
}

func getUser(c *gin.Context) {
	id := c.Param("id")
	user, ok := userMap.Get(id)
	if !ok {
		c.JSON(404, "not-found")
		return
	}
	c.JSON(200, user)
}

func createUser(c *gin.Context) {
	val := atomic.AddUint64(&idCounter, 1)
	user := new(User)
	if e := c.BindJSON(user); e != nil {
		log.Println(e)
		c.JSON(401, "bad-request")
		return
	}
	user.Id = fmt.Sprintf("u%v", val)
	userMap.Set(user.Id, user)
	c.JSON(200, user)
}

func updateUser(c *gin.Context) {
	id := c.Param("id")
	u, ok := userMap.Get(id)
	if !ok {
		c.JSON(404, "not-found")
		return
	}
	user := u.(*User)
	c.BindJSON(user)
	user.Id = id
	c.JSON(200, user)
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")
	_, ok := userMap.Get(id)
	if !ok {
		c.JSON(404, "not-found")
		return
	}
	userMap.Remove(id)
	c.JSON(200, "ok")
}

func setupJSON() {
	extra.SetNamingStrategy(func(name string) string {
		//make json key names lower camel-case
		n := name[0:1]
		n = strings.ToLower(n)
		if len(name) > 1 {
			n += name[1:]
		}
		return n
	})
}

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.json", "path/to/config.json")
	flag.Parse()

	setupJSON()

	config.Init(configPath)
	services.Init()

	// router := gin.Default()
	// api := router.Group("/api")
	// api.GET("/users/:id", getUser)
	// api.POST("/users", createUser)
	// api.PUT("/users/:id", updateUser)
	// api.DELETE("/users/:id", deleteUser)
	// router.Run()
}
