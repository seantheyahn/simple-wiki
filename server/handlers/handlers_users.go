package handlers

import (
	"regexp"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/seantheyahn/simple-wiki/services"
)

type projectUser struct {
	services.User
	CanWrite bool
}

var regexUsername = regexp.MustCompile(`^[A-Za-z_]{1}[A-Za-z_0-9]{3,}$`)

func checkPasswordFormat(password string) bool {
	check := func(patterns ...string) bool {
		for _, pattern := range patterns {
			if m, _ := regexp.MatchString(pattern, password); !m {
				return false
			}
		}
		return true
	}
	return check("[a-z]+", "[A-Z]+", "[0-9]+", ".{8,}")
}

func usersIndex(c *gin.Context) {
	var users []*services.User
	var err error
	user := getUser(c)
	if user.Admin {
		users, err = services.LoadAllUsers()
	} else {
		u, err := services.LoadUser(user.ID)
		if checkError(c, err) {
			return
		}
		users = []*services.User{u}
	}
	if checkError(c, err) {
		return
	}
	sort.Slice(users, func(i int, j int) bool {
		return users[i].CreatedAt.Unix() < users[j].CreatedAt.Unix()
	})
	c.HTML(200, "users.html", gin.H{
		"User":  user,
		"Users": users,
	})
}

func projectUsersIndex(c *gin.Context) {
	user := getUser(c)
	id, err := strconv.Atoi(c.Param("id"))
	if checkError(c, err) {
		return
	}
	project, err := services.LoadProject(id)
	if checkError(c, err) {
		return
	}
	if !user.Admin {
		role, err := services.LoadRole(id, user.ID)
		if checkError(c, err) {
			return
		}
		if !role.CanWrite {
			c.String(403, "permission-denied")
			return
		}
	}
	allUsers, err := services.LoadAllUsers()
	if checkError(c, err) {
		return
	}
	roles, err := services.LoadProjectRoles(id)
	if checkError(c, err) {
		return
	}
	userMap := make(map[string]*services.User)
	for _, u := range allUsers {
		userMap[u.ID] = u
	}
	projectUsers := make([]*projectUser, 0, len(roles))
	for _, role := range roles {
		pu := new(projectUser)
		pu.User = *userMap[role.UserID]
		pu.CanWrite = role.CanWrite
		projectUsers = append(projectUsers, pu)
	}

	sort.Slice(allUsers, func(i int, j int) bool {
		return allUsers[i].CreatedAt.Unix() < allUsers[j].CreatedAt.Unix()
	})
	sort.Slice(projectUsers, func(i int, j int) bool {
		return projectUsers[i].CreatedAt.Unix() < projectUsers[j].CreatedAt.Unix()
	})
	c.HTML(200, "project-users.html", gin.H{
		"User":         user,
		"Project":      project,
		"AllUsers":     allUsers,
		"ProjectUsers": projectUsers,
	})
}

func updateRole(c *gin.Context) {
	type form struct {
		UserID    string `form:"userID" binding:"required"`
		ProjectID int    `form:"projectID" binding:"required"`
		ReadOnly  string `form:"readonly"`
	}
	action := c.Param("action")
	user := getUser(c)
	f := new(form)
	if checkError(c, c.Bind(f)) {
		return
	}
	if !user.Admin {
		role, err := services.LoadRole(f.ProjectID, user.ID)
		if checkError(c, err) {
			return
		}
		if !role.CanWrite {
			c.String(403, "permission-denied")
			return
		}
	}
	var err error
	switch action {
	case "add":
		err = services.AddRole(f.ProjectID, f.UserID, f.ReadOnly != "on")
	case "remove":
		err = services.DeleteRole(f.ProjectID, f.UserID)
	case "edit":
		err = services.UpdateRole(f.ProjectID, f.UserID, f.ReadOnly != "on")
	}

	if checkError(c, err) {
		return
	}
	c.Redirect(302, "/projects/users/view/"+strconv.Itoa(f.ProjectID))
}

func createUser(c *gin.Context) {
	type form struct {
		Username string `form:"username" binding:"required"`
		Password string `form:"password" binding:"required"`
		Admin    string `form:"admin"`
	}
	f := new(form)
	if checkError(c, c.Bind(f)) {
		return
	}
	if !regexUsername.Match([]byte(f.Username)) || !checkPasswordFormat(f.Password) {
		c.String(400, "bad-request")
		return
	}
	_, err := services.CreateUser(f.Username, f.Password, f.Admin == "on", false)
	if checkError(c, err) {
		return
	}
	c.Redirect(302, "/users")
}

func editUser(c *gin.Context) {
	type form struct {
		ID       string `form:"id" binding:"required"`
		Username string `form:"username" binding:"required"`
		Admin    string `form:"admin"`
	}
	f := new(form)
	if checkError(c, c.Bind(f)) {
		return
	}
	user := getUser(c)
	if !user.Admin && user.ID != f.ID {
		c.String(403, "permission-denied")
		return
	}
	if !regexUsername.Match([]byte(f.Username)) {
		c.String(400, "bad-request")
		return
	}
	err := services.UpdateUser(f.ID, f.Username, f.Admin == "on")
	if checkError(c, err) {
		return
	}
	c.Redirect(302, "/users")
}

func changeUserPassword(c *gin.Context) {
	type form struct {
		ID       string `form:"id" binding:"required"`
		Password string `form:"password" binding:"required"`
	}
	f := new(form)
	if checkError(c, c.Bind(f)) {
		return
	}
	user := getUser(c)
	if !user.Admin && user.ID != f.ID {
		c.String(403, "permission-denied")
		return
	}
	if !checkPasswordFormat(f.Password) {
		c.String(400, "bad-request")
		return
	}
	if checkError(c, services.ChangeUserPassword(f.ID, f.Password)) {
		return
	}
	c.Redirect(302, "/users")
}

func deleteUser(c *gin.Context) {
	type form struct {
		ID string `form:"id" binding:"required"`
	}
	f := new(form)
	if checkError(c, c.Bind(f)) {
		return
	}
	err := services.DeleteUser(f.ID)
	if checkError(c, err) {
		return
	}
	c.Redirect(302, "/users")
}
