package handlers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/seantheyahn/simple-wiki/services"
	csrf "github.com/utrack/gin-csrf"
)

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
			services.Logger.Error(err)
			c.String(500, "internal-server-error")
			return
		}

		c.Redirect(302, "/")
		return
	}

	c.HTML(200, "login.html", data)
}

func logout(c *gin.Context) {
	s := sessions.Default(c)
	s.Delete("user")
	if s.Save() != nil {
		c.String(500, "internal-server-error")
		return
	}
	c.Redirect(302, "/")
}

func authMiddleware(c *gin.Context) {
	s := sessions.Default(c)
	u := s.Get("user")
	if u == nil {
		c.String(401, "unauthorized")
		c.Abort()
		return
	}
	c.Set("user", u)
	c.Next()
}

func getUser(c *gin.Context) *services.User {
	u, _ := c.Get("user")
	return u.(*services.User)
}
