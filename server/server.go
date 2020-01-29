package server

import (
	"html/template"
	"log"
	"path/filepath"
	"time"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/seantheyahn/simple-wiki/config"
	"github.com/seantheyahn/simple-wiki/server/handlers"
	"github.com/seantheyahn/simple-wiki/services"
	limit "github.com/yangxikun/gin-limit-by-key"
	"golang.org/x/time/rate"
)

func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	layouts, err := filepath.Glob(templatesDir + "/layouts/*.html")
	if err != nil {
		panic(err.Error())
	}

	includes, err := filepath.Glob(templatesDir + "/includes/*.html")
	if err != nil {
		panic(err.Error())
	}

	// Generate our templates map from our layouts/ and includes/ directories
	for _, include := range includes {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, include)
		r.AddFromFilesFuncs(filepath.Base(include), template.FuncMap{
			"htmlSafe": func(html string) template.HTML {
				return template.HTML(html)
			},
		}, files...)
	}
	return r
}

//Run runs the http server blocking mode
func Run() {
	if config.Instance.Server.DevelopmentMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.HTMLRender = loadTemplates("./views")

	//logger middleware
	router.Use(services.PanicRecoveryMiddleware())
	router.Use(services.LoggerMiddleware())

	//sessions middleware
	router.Use(sessions.Sessions("user", cookie.NewStore([]byte(config.Instance.Server.CookieSecret))))

	//rate limiter middleware
	router.Use(limit.NewRateLimiter(func(c *gin.Context) string {
		return c.ClientIP() // limit rate by client ip
	}, func(c *gin.Context) (*rate.Limiter, time.Duration) {
		interval := time.Duration(int(time.Second) / config.Instance.Server.IPRateLimiter.ReqPerSecond)
		burst := config.Instance.Server.IPRateLimiter.ReqBurst
		return rate.NewLimiter(rate.Every(interval), burst), time.Minute
	}, func(c *gin.Context) {
		c.String(429, "too-many-requests")
		c.Abort()
	}))

	//add routes
	handlers.AddHandlers(router)

	services.Logger.Infof("listening on %v", config.Instance.Server.ListenAddress)
	log.Fatal(router.Run(config.Instance.Server.ListenAddress))
}
