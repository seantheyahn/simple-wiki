package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/seantheyahn/simple-wiki/services"
)

func createDocument(c *gin.Context) {
	type form struct {
		ProjectID int    `form:"projectID" binding:"required"`
		Title     string `form:"title" binding:"required"`
		Body      string `form:"body" binding:"required"`
		SortOrder int    `form:"sortOrder"`
	}
	user := getUser(c)
	doc := new(form)
	if checkError(c, c.Bind(doc)) {
		return
	}
	if doc.Title == "" {
		c.String(400, "bad-request")
		return
	}

	if !user.Admin {
		pu, err := services.LoadRole(doc.ProjectID, user.ID)
		if checkError(c, err) {
			return
		}
		if !pu.CanWrite {
			c.String(403, "permission-denied")
			return
		}
	}
	_, err := services.CreateDocument(doc.ProjectID, doc.Title, doc.Body, doc.SortOrder)
	if checkError(c, err) {
		return
	}
	c.Redirect(302, fmt.Sprintf("/projects/view/%d", doc.ProjectID))
}

func editDocument(c *gin.Context) {
	type form struct {
		ProjectID  int    `form:"projectID" binding:"required"`
		DocumentID int    `form:"documentID" binding:"required"`
		Title      string `form:"title" binding:"required"`
		Body       string `form:"body" binding:"required"`
		SortOrder  int    `form:"sortOrder"`
	}
	user := getUser(c)
	doc := new(form)
	if checkError(c, c.Bind(doc)) {
		return
	}
	if doc.Title == "" {
		c.String(400, "bad-request")
		return
	}
	//check if the document is related to provided the project id
	d, err := services.LoadDocument(doc.DocumentID)
	if checkError(c, err) {
		return
	}
	if d.ProjectID != doc.ProjectID {
		c.String(400, "bad-request")
		return
	}

	if !user.Admin {
		pu, err := services.LoadRole(doc.ProjectID, user.ID)
		if checkError(c, err) {
			return
		}
		if !pu.CanWrite {
			c.String(403, "permission-denied")
			return
		}
	}
	if checkError(c, services.UpdateDocument(doc.DocumentID, doc.Title, doc.Body, doc.SortOrder)) {
		return
	}
	c.Redirect(302, fmt.Sprintf("/projects/view/%d", doc.ProjectID))
}

func deleteDocument(c *gin.Context) {
	type form struct {
		ProjectID  int `form:"projectID" binding:"required"`
		DocumentID int `form:"documentID" binding:"required"`
	}
	user := getUser(c)
	f := new(form)
	if checkError(c, c.Bind(f)) {
		return
	}
	//check if the document is related to provided the project id
	d, err := services.LoadDocument(f.DocumentID)
	if checkError(c, err) {
		return
	}
	if d.ProjectID != f.ProjectID {
		c.String(400, "bad-request")
		return
	}

	if !user.Admin {
		pu, err := services.LoadRole(f.ProjectID, user.ID)
		if checkError(c, err) {
			return
		}
		if !pu.CanWrite {
			c.String(403, "permission-denied")
			return
		}
	}
	if checkError(c, services.DeleteDocument(f.DocumentID)) {
		return
	}
	c.Redirect(302, fmt.Sprintf("/projects/view/%d", f.ProjectID))
}
