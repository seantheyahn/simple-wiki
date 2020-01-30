package handlers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/seantheyahn/simple-wiki/config"
	"github.com/seantheyahn/simple-wiki/services"
	csrf "github.com/utrack/gin-csrf"
)

func index(c *gin.Context) {
	s := sessions.Default(c)
	u := s.Get("user")
	if u == nil {
		c.Redirect(302, "/auth/login")
		return
	}

	c.Redirect(302, "/projects")
}

//AddHandlers adds routes to the provided router
func AddHandlers(router *gin.Engine) {
	router.GET("/", index)
	{
		r := router.Group("/auth")
		//CSRF middleware
		r.Use(csrf.Middleware(csrf.Options{
			Secret: config.Instance.Server.CSRFSecret,
			ErrorFunc: func(c *gin.Context) {
				c.String(400, "CSRF token mismatch")
				c.Abort()
			},
		}))
		r.GET("/login", login)
		r.POST("/login", login)
		r.GET("/logout", logout)
	}
	privateRouter := router.Group("/")
	privateRouter.Use(authMiddleware)
	{
		r := privateRouter.Group("/projects")
		r.GET("/", projectsIndex)
		r.GET("/users/view/:id", projectUsersIndex)
		r.POST("/users/update/:action", updateRole)
		r.GET("/view/:id", viewProject)
		r.POST("/edit/:id", editProject)
		r.POST("/delete/:id", deleteProject)
		r.POST("/create", createProject)
	}
	{
		r := privateRouter.Group("/documents")
		r.POST("/edit", editDocument)
		r.POST("/delete", deleteDocument)
		r.POST("/create", createDocument)
	}
	{
		r := privateRouter.Group("/users")
		r.GET("/", usersIndex)
		r.POST("/edit", editUser)
		r.POST("/delete", adminCheckMiddleware, deleteUser)
		r.POST("/create", adminCheckMiddleware, createUser)
		r.POST("/password", changeUserPassword)
	}
}

func checkError(c *gin.Context, err error) bool {
	if err != nil {
		services.Logger.Error(err)
		c.String(500, "internal-server-error")
		c.Abort()
		return true
	}
	return false
}
