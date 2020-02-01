package services

import (
	"database/sql"
	"encoding/gob"
)

func checkPanic(err error) {
	if err != nil {
		panic(err)
	}
}

//DB database instance
var DB *sql.DB

//Init --
func Init() {
	initLogger()
	Logger.Info("initializing services")
	DB = initDB()
	migrateDB()
	initUsers()

	gob.Register(new(User))
}

//Finalize --
func Finalize() {
	finalizeLogger()
}
