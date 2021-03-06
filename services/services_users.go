package services

import (
	"database/sql"
	"errors"
	"time"

	"github.com/seantheyahn/simple-wiki/config"

	"github.com/google/uuid"
	"github.com/jackc/pgx"
)

//User user model
type User struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	ID        string
	Username  string
	Admin     bool
}

//ErrUsernameAlreadyExists --
var ErrUsernameAlreadyExists = errors.New("username already exists")

//ErrAuthenticationFailed --
var ErrAuthenticationFailed = errors.New("authentication failed")

//ErrUserNotFound --
var ErrUserNotFound = errors.New("user not found")

func initUsers() {
	//create root user if not exists
	u := config.Instance.RootUser
	user, err := CreateUser(u.Username, u.Password, true, true)
	//error code ErrUsernameAlreadyExists also works for id, which is <root> for root user
	if err != nil && err != ErrUsernameAlreadyExists {
		panic(err)
	}
	if user != nil {
		Logger.Info("created root user with username: ", user.Username)
	}
}

//CreateUser --
func CreateUser(username string, password string, admin bool, root bool) (*User, error) {
	user := &User{
		Username: username,
		Admin:    admin,
	}
	if root {
		user.ID = "<root>"
	} else {
		user.ID = uuid.New().String()
	}
	passwordHash, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	row := DB.QueryRow("insert into users (id, username, password_hash, admin) values($1,$2,$3,$4) returning created_at, updated_at", user.ID, user.Username, passwordHash, user.Admin)

	if err := row.Scan(&user.CreatedAt, &user.UpdatedAt); err != nil {
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.SQLState() == "23505" {
			//duplicate key value
			return nil, ErrUsernameAlreadyExists
		}

		return nil, err
	}
	return user, nil
}

//LoadAllUsers --
func LoadAllUsers() (result []*User, err error) {
	result = make([]*User, 0)
	rows, err := DB.Query("select id, created_at, updated_at, username, admin from users")
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		u := new(User)
		err = rows.Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt, &u.Username, &u.Admin)
		if err != nil {
			return
		}
		result = append(result, u)
	}
	return
}

//LoadUser --
func LoadUser(id string) (*User, error) {
	user := new(User)
	user.ID = id
	err := DB.QueryRow("select created_at, updated_at, username, admin from users where id=$1", id).Scan(&user.CreatedAt, &user.UpdatedAt, &user.Username, &user.Admin)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

//UserFindIDByUserName returns the id of the user
func UserFindIDByUserName(username string) (id string, err error) {
	err = DB.QueryRow("select id from users where username=$1", username).Scan(&id)
	return
}

//AuthenticateUser --
func AuthenticateUser(username string, password string) (user *User, err error) {
	var passwordHash string
	var id string
	err = DB.QueryRow("select id, password_hash from users where username=$1", username).Scan(&id, &passwordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrUserNotFound
			return
		}
		return
	}
	if !CheckPasswordHash(password, passwordHash) {
		err = ErrAuthenticationFailed
		return
	}
	user, err = LoadUser(id)
	return
}

//DeleteUser --
func DeleteUser(id string) error {
	res, err := DB.Exec("delete from users where id=$1", id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if n == 0 {
		err = ErrUserNotFound
	}
	return err
}

//UpdateUser fields: Username, Admin
func UpdateUser(id string, username string, admin bool) error {
	res, err := DB.Exec("update users set username=$1, admin=$2 where id=$3", username, admin, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if n == 0 {
		err = ErrUserNotFound
	}
	return err
}

//ChangeUserPassword --
func ChangeUserPassword(id string, newPassword string) error {
	hash, err := HashPassword(newPassword)
	if err != nil {
		return err
	}
	res, err := DB.Exec("update users set password_hash=$1 where id=$2", hash, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if n == 0 {
		err = ErrUserNotFound
	}
	return err
}
