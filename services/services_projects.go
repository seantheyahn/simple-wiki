package services

import (
	"time"
)

//Project project model
type Project struct {
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ID          int
	Title       string
	Description string
}

//ProjectUser project user model
type ProjectUser struct {
	ProjectID int
	UserID    string
	CanWrite  bool
}

//CreateProject --
func CreateProject(title string, description string, user *User) (project *Project, err error) {
	project = &Project{
		Title:       title,
		Description: description,
	}

	tx, err := DB.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()
	err = tx.QueryRow("insert into projects (title, description) values($1,$2) returning created_at, updated_at, id", title, description).Scan(&project.CreatedAt, &project.UpdatedAt, &project.ID)
	if err != nil {
		return
	}
	_, err = tx.Exec("insert into projects__users (project_id, user_id, can_write) values($1,$2,$3)", project.ID, user.ID, true)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

//LoadProject --
func LoadProject(id int) (*Project, error) {
	project := new(Project)
	project.ID = id
	err := DB.QueryRow("select created_at, updated_at, title, description from projects where id=$1", id).Scan(&project.CreatedAt, &project.UpdatedAt, &project.Title, &project.Description)
	if err != nil {
		return nil, err
	}
	return project, nil
}

//LoadProjectIDsForUser --
func LoadProjectIDsForUser(userID string) ([]int, error) {
	rows, err := DB.Query("select project_id from projects__users where user_id=$1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]int, 0)
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		result = append(result, id)
	}
	return result, nil
}

//LoadAllProjectIDs --
func LoadAllProjectIDs() ([]int, error) {
	rows, err := DB.Query("select id from projects")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]int, 0)
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		result = append(result, id)
	}
	return result, nil
}

//LoadProjectUser --
func LoadProjectUser(projectID int, userID string) (*ProjectUser, error) {
	pu := new(ProjectUser)
	pu.ProjectID = projectID
	pu.UserID = userID
	row := DB.QueryRow("select can_write from projects__users where project_id=$1 and user_id=$2", projectID, userID)
	err := row.Scan(&pu.CanWrite)
	return pu, err
}

//LoadProjectUsers --
func LoadProjectUsers(projectID int) ([]*ProjectUser, error) {
	rows, err := DB.Query("select project_id, user_id, can_write from projects__users where project_id=$1", projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]*ProjectUser, 0)
	for rows.Next() {
		pu := new(ProjectUser)
		err = rows.Scan(&pu.ProjectID, &pu.UserID, &pu.CanWrite)
		if err != nil {
			return nil, err
		}
		result = append(result, pu)
	}
	return result, nil
}

//UpdateProject --
func UpdateProject(id int, title string, description string) error {
	_, err := DB.Exec("update projects set title=$1, description=$2 where id=$3", title, description, id)
	return err
}

//DeleteProject --
func DeleteProject(id int) error {
	_, err := DB.Exec("delete from projects where id=$1", id)
	return err
}
