package server

import (
	"github.com/gin-gonic/gin"
	"log"
	"sean/wiki/config"
)

//Run runs the http server blocking mode
func Run() {
	router := gin.Default()
	router.LoadHTMLGlob("views/*.html")

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	log.Fatal(router.Run(config.Instance.Server.ListenAddress))
}
