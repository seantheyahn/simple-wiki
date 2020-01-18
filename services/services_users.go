package services

import (
	"github.com/google/uuid"
)

//User user model
type User struct {
	ID   string
	Name string
}

func initUsers() {

}

//CreateUser --
func CreateUser(name string) (*User, error) {
	user := &User{
		ID:   uuid.New().String(),
		Name: name,
	}
	if _, err := DB.Exec("insert into users (id, name) values($1, $2)", user.ID, user.Name); err != nil {
		return nil, err
	}
	return user, nil
}

//LoadUser -z
func LoadUser(id string) (*User, error) {
	return nil, nil
}

//DeleteUser --
func DeleteUser(id string) error {
	return nil
}

//UpdateUser --
func UpdateUser(user *User) (*User, error) {
	return nil, nil
}
