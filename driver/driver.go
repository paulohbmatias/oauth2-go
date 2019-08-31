package driver

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var db *sql.DB

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ConnectDB() *sql.DB {

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s "+ "password=%s dbname=%s sslmode=disable", os.Getenv("HOST"),
		os.Getenv("PORT_DB"), os.Getenv("USER_DB"), os.Getenv("PASSWORD"), os.Getenv("DBNAME"))

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	logFatal(err)

	err = db.Ping()
	logFatal(err)

	return db
}