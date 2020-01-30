package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/seantheyahn/simple-wiki/services"
)

type renderedProject struct {
	services.Project
	CanEdit bool
}

func projectsIndex(c *gin.Context) {
	data := new(struct {
		User         *services.User
		UserProjects []*renderedProject
		AllProjects  []*services.Project
	})
	data.User = getUser(c)
	ids, err := services.LoadProjectIDsForUser(data.User.ID)
	if checkError(c, err) {
		return
	}
	writeMap := make(map[int]bool)
	if !data.User.Admin {
		roles, err := services.LoadUserRoles(data.User.ID)
		if checkError(c, err) {
			return
		}
		for _, r := range roles {
			writeMap[r.ProjectID] = r.CanWrite
		}
	}

	data.UserProjects = make([]*renderedProject, 0, len(ids))
	for _, id := range ids {
		p, err := services.LoadProject(id)
		if checkError(c, err) {
			return
		}
		rp := new(renderedProject)
		rp.Project = *p
		rp.CanEdit = data.User.Admin || writeMap[p.ID]
		data.UserProjects = append(data.UserProjects, rp)
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
		Description string `form:"description"`
	}
	user := getUser(c)
	project := new(form)

	c.Bind(project)
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
		Description string `form:"description"`
	}
	user := getUser(c)
	project := new(form)
	id, err := strconv.Atoi(c.Param("id"))
	if checkError(c, err) {
		return
	}
	if !user.Admin {
		pu, err := services.LoadRole(id, user.ID)
		if checkError(c, err) {
			return
		}
		if !pu.CanWrite {
			c.String(403, "permission-denied")
			return
		}
	}
	p, err := services.LoadProject(id)
	if checkError(c, err) {
		return
	}
	project.Title = p.Title
	project.Description = p.Description

	c.Bind(project)
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
	user := getUser(c)
	id, err := strconv.Atoi(c.Param("id"))
	if checkError(c, err) {
		return
	}
	if !user.Admin {
		pu, err := services.LoadRole(id, user.ID)
		if checkError(c, err) {
			return
		}
		if !pu.CanWrite {
			c.String(403, "permission-denied")
			return
		}
	}
	if checkError(c, services.DeleteProject(id)) {
		return
	}
	c.Redirect(302, "/projects")
}

func viewProject(c *gin.Context) {
	data := new(struct {
		User      *services.User
		Project   *services.Project
		Documents []*services.Document
		Role      *services.Role
	})
	id, err := strconv.Atoi(c.Param("id"))
	if checkError(c, err) {
		return
	}
	data.Project, err = services.LoadProject(id)
	if checkError(c, err) {
		return
	}
	data.User = getUser(c)
	if data.User.Admin {
		data.Role = &services.Role{ProjectID: data.Project.ID, UserID: data.User.ID, CanWrite: true}
	} else {
		data.Role, err = services.LoadRole(id, data.User.ID)
		if checkError(c, err) {
			return
		}
	}
	data.Documents, err = services.LoadDocuments(id)
	if checkError(c, err) {
		return
	}
	c.HTML(200, "documents.html", data)
}
