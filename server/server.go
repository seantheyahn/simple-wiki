package server

import (
	"encoding/gob"
	"log"
	"path/filepath"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/seantheyahn/simple-wiki/config"
	"github.com/seantheyahn/simple-wiki/services"
	csrf "github.com/utrack/gin-csrf"
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
		log.Println(r)
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

	c.HTML(200, "index.html", u)
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
	router := gin.Default()
	router.HTMLRender = loadTemplates("./views")
	router.Use(sessions.Sessions("user", cookie.NewStore([]byte(config.Instance.Server.CookieSecret))))
	router.Use(csrf.Middleware(csrf.Options{
		Secret: config.Instance.Server.CSRFSecret,
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	}))

	gob.Register(&services.User{})

	router.GET("/", index)
	router.GET("/login", login)
	router.POST("/login", login)
	router.GET("/logout", logout)

	log.Fatal(router.Run(config.Instance.Server.ListenAddress))
}
