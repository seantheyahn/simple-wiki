package services

import (
	"database/sql"
	"io/ioutil"
	"log"
	"net/url"
	"path/filepath"
	"regexp"
	"sean/wiki/config"
	"strconv"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
)

func initDB() *sql.DB {
	dbConf := config.Instance.Db
	uri := dbConf.ConnectionURI
	log.Println("initializing database", uri)
	parsedURL, err := url.Parse(uri)
	checkPanic(err)
	db, err := sql.Open("pgx", parsedURL.String())
	checkPanic(err)
	db.SetMaxOpenConns(dbConf.MaxOpenConns)
	db.SetMaxIdleConns(dbConf.MaxIdleConns)

	var dbVersion string
	checkPanic(db.QueryRow("select version()").Scan(&dbVersion))
	log.Println("database version:", dbVersion)
	return db
}

func migrateDB() {
	const migrationsRoot = "migrations"
	reDir := regexp.MustCompile("^[0-9]+\\.([a-zA-Z_]+)$")
	reFile := regexp.MustCompile("^([0-9]+)\\.sql$")
	_dirs, err := ioutil.ReadDir(migrationsRoot)
	checkPanic(err)
	conn, err := stdlib.AcquireConn(DB)
	checkPanic(err)
	defer stdlib.ReleaseConn(DB, conn)
	log.Println("performing db migrations")
	_, err = conn.Exec(`create table if not exists _migrations (
			key varchar primary key,
			version int not null default 0
		);`)
	checkPanic(err)

	for _, entry := range _dirs {
		if entry.IsDir() && reDir.MatchString(entry.Name()) {
			key := reDir.FindStringSubmatch(entry.Name())[1]
			_dirPath := filepath.Join(migrationsRoot, entry.Name())
			_files, err := ioutil.ReadDir(_dirPath)
			checkPanic(err)
			for _, _file := range _files {
				if !_file.IsDir() && reFile.MatchString(_file.Name()) {
					_filePath := filepath.Join(_dirPath, _file.Name())
					v, err := strconv.Atoi(reFile.FindStringSubmatch(_file.Name())[1])
					checkPanic(err)
					currentV := 0
					err = conn.QueryRow(`select version from _migrations where key=$1`, key).Scan(&currentV)
					if err != pgx.ErrNoRows {
						checkPanic(err)
					}
					if v <= currentV {
						continue
					}
					q, err := ioutil.ReadFile(_filePath)
					checkPanic(err)
					log.Printf("migrating '%s' from version %v to %v\n", key, currentV, v)
					tx, err := conn.Begin()
					checkPanic(err)
					_, err = tx.Exec(string(q))
					checkPanic(err)
					_, err = tx.Exec("insert into _migrations (key, version) values ($1, $2) on conflict(key) do update set version=EXCLUDED.version;", key, v)
					checkPanic(err)
					checkPanic(tx.Commit())
				}
			}
		}
	}
	log.Println("finished db migrations")
}
