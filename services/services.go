package services

import (
	"database/sql"
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

	test() //TODO remove
}

//Finalize --
func Finalize() {
	finalizeLogger()
}

func test() {
}
