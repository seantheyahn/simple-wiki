package server

import (
	"encoding/gob"
	"log"
	"path/filepath"
	"time"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/seantheyahn/simple-wiki/config"
	"github.com/seantheyahn/simple-wiki/services"
	csrf "github.com/utrack/gin-csrf"
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
		r.AddFromFiles(filepath.Base(include), files...)
	}
	return r
}

func index(c *gin.Context) {
	s := sessions.Default(c)
	u := s.Get("user")
	if u == nil {
		c.Redirect(302, "/login")
		return
	}

	data := struct {
		User interface{}
	}{u}

	c.HTML(200, "index.html", data)
}

func login(c *gin.Context) {
	s := sessions.Default(c)
	u := s.Get("user")
	if u != nil {
		c.Redirect(302, "/")
		return
	}

	data := new(struct {
		Username string `form:"username" binding:"required"`
		Password string `form:"password" binding:"required"`
		Err      string
		CSRF     string
		User     *services.User
	})

	data.CSRF = csrf.GetToken(c)

	if c.Request.Method == "POST" {
		if err := c.Bind(&data); err != nil {
			c.String(400, "bad-request")
			return
		}
		user, err := services.AuthenticateUser(data.Username, data.Password)
		if err != nil {
			switch err {
			case services.ErrAuthenticationFailed, services.ErrUserNotFound:
				data.Err = "Login Failed!"
				c.HTML(200, "login.html", data)
				return
			default:
				c.Status(400)
				return
			}
		}
		s.Set("user", user)
		if err = s.Save(); err != nil {
			log.Println(err)
			c.String(500, "internal-server-error")
			return
		}

		c.Redirect(302, "/")
		return
	}

	c.HTML(200, "login.html", data)
}

func logout(c *gin.Context) {
	log.Println("logging out")
	s := sessions.Default(c)
	s.Delete("user")
	if s.Save() != nil {
		c.String(500, "internal-server-error")
		return
	}
	c.Redirect(302, "/")
}

//Run runs the http server blocking mode
func Run() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.HTMLRender = loadTemplates("./views")
	router.Use(gin.Recovery())

	//logger middleware
	router.Use(ginzap.Ginzap(services.LoggerCore, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(services.LoggerCore, true))

	//sessions middleware
	router.Use(sessions.Sessions("user", cookie.NewStore([]byte(config.Instance.Server.CookieSecret))))

	//CSRF middleware
	router.Use(csrf.Middleware(csrf.Options{
		Secret: config.Instance.Server.CSRFSecret,
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	}))

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

	//register types used in sessions to gob
	gob.Register(&services.User{})

	//routes
	router.GET("/", index)
	router.GET("/login", login)
	router.POST("/login", login)
	router.GET("/logout", logout)

	services.Logger.Infof("listening on %v", config.Instance.Server.ListenAddress)
	log.Fatal(router.Run(config.Instance.Server.ListenAddress))
}
