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
	data := new(struct {
		User    *services.User
		Project *form
		Err     string
	})
	data.User = getUser(c)
	data.Project = new(form)

	if c.Request.Method == "POST" {
		c.Bind(&data.Project)
		if data.Project.Title == "" {
			data.Err = "Project title cannot be empty"
		} else {
			_, err := services.CreateProject(data.Project.Title, data.Project.Description, data.User)
			if checkError(c, err) {
				return
			}
			c.Redirect(302, "/projects")
			return
		}
	}

	c.HTML(200, "create-project.html", data)
}

func editProject(c *gin.Context) {
	type form struct {
		Title       string `form:"title" binding:"required"`
		Description string `form:"description" binding:"required"`
	}
	data := new(struct {
		User    *services.User
		Project *form
		Err     string
	})
	data.User = getUser(c)
	data.Project = new(form)
	id, err := strconv.Atoi(c.Param("id"))
	if checkError(c, err) {
		return
	}
	pu, err := services.LoadProjectUser(id, data.User.ID)
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
	data.Project.Title = p.Title
	data.Project.Description = p.Description

	if c.Request.Method == "POST" {
		c.Bind(&data.Project)
		if data.Project.Title == "" {
			data.Err = "Project title cannot be empty"
		} else {
			err := services.UpdateProject(id, data.Project.Title, data.Project.Description)
			if checkError(c, err) {
				return
			}
			c.Redirect(302, "/projects")
			return
		}
	}

	c.HTML(200, "edit-project.html", data)
}
