package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/seantheyahn/simple-wiki/services"
	"strconv"
)

func projectsIndex(c *gin.Context) {
	data := new(struct {
		User         *services.User
		UserProjects []*services.Project
		AllProjects  []*services.Project
	})
	data.User = getUser(c)
	ids, err := services.LoadProjectIDsForUser(data.User.ID)
	if checkError(c, err) {
		return
	}
	data.UserProjects = make([]*services.Project, 0, len(ids))
	for _, id := range ids {
		p, err := services.LoadProject(id)
		if checkError(c, err) {
			return
		}
		data.UserProjects = append(data.UserProjects, p)
	}
	if data.User.Admin {
		ids, err = services.LoadAllProjectIDs()
		data.AllProjects = make([]*services.Project, 0, len(ids))
		for _, id := range ids {
			p, err := services.LoadProject(id)
			if checkError(c, err) {
				return
			}
			data.AllProjects = append(data.AllProjects, p)
		}
	}

	c.HTML(200, "projects.html", data)
}

func createProject(c *gin.Context) {
	type form struct {
		Title       string `form:"title" binding:"required"`
		Description string `form:"description" binding:"required"`
	}
	user := getUser(c)
	project := new(form)

	c.Bind(&project)
	if project.Title == "" {
		c.String(400, "bad-request")
		return
	}

	_, err := services.CreateProject(project.Title, project.Description, user)
	if checkError(c, err) {
		return
	}
	c.Redirect(302, "/projects")
}

func editProject(c *gin.Context) {
	type form struct {
		Title       string `form:"title" binding:"required"`
		Description string `form:"description" binding:"required"`
	}
	user := getUser(c)
	project := new(form)
	id, err := strconv.Atoi(c.Param("id"))
	if checkError(c, err) {
		return
	}
	pu, err := services.LoadProjectUser(id, user.ID)
	if checkError(c, err) {
		return
	}
	if !pu.CanWrite {
		c.String(403, "permission-denied")
		return
	}
	p, err := services.LoadProject(id)
	if checkError(c, err) {
		return
	}
	project.Title = p.Title
	project.Description = p.Description

	c.Bind(&project)
	if project.Title == "" {
		c.String(400, "bad-request")
		return
	}

	err = services.UpdateProject(id, project.Title, project.Description)
	if checkError(c, err) {
		return
	}
	c.Redirect(302, "/projects")
}

func deleteProject(c *gin.Context) {
	type form struct {
		Title       string `form:"title" binding:"required"`
		Description string `form:"description" binding:"required"`
	}
	user := getUser(c)
	id, err := strconv.Atoi(c.Param("id"))
	if checkError(c, err) {
		return
	}
	pu, err := services.LoadProjectUser(id, user.ID)
	if checkError(c, err) {
		return
	}
	if !pu.CanWrite {
		c.String(403, "permission-denied")
		return
	}
	if checkError(c, services.DeleteProject(id)) {
		return
	}
	c.Redirect(302, "/projects")
}
