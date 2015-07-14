package tdb

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

import _ "github.com/go-sql-driver/mysql"

var (
	conn *sql.DB
)

func init() {
	var err error
	if conn, err = sql.Open("mysql", fmt.Sprintf("%s:%s@/%s",
		os.Getenv("dbuser"), os.Getenv("dbpwd"), os.Getenv("dbdbs"))); err != nil {
		log.Fatal(err)
	}
}
