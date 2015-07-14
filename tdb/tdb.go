package tdb

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var conn *sql.DB

type base struct {
	db    *sql.DB
	table string
}

func init() {
	var err error
	if conn, err = sql.Open("mysql", fmt.Sprintf("%s:%s@/%s",
		os.Getenv("dbuser"), os.Getenv("dbpwd"), os.Getenv("dbdbs"))); err != nil {
		log.Fatal(err)
	}
}
