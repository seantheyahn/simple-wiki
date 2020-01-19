package services

import (
	"database/sql"
	"log"
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
	log.Println("initializing services")
	DB = initDB()
	initRedis()
	migrateDB()
	initUsers()

	test() //TODO remove
}

func test() {

}
