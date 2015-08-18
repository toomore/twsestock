package tdb

import (
	"database/sql"
	"fmt"
	"log"
)

type filterinfo base

var insertFilterinfoSQL *sql.Stmt

func NewFilterinfoDB() *filterinfo {
	table := "filterinfo"
	if insertFilterinfoSQL, err = conn.Prepare(fmt.Sprintf("Insert into %s(no, description) Values(?,?) ON DUPLICATE KEY UPDATE no=?, description=?", table)); err != nil {
		log.Fatal(err)
	}

	return &filterinfo{
		table: table,
	}
}

func (filterinfo) InsertFilterinfo(no uint64, desc string) (sql.Result, error) {
	return insertFilterinfoSQL.Exec(no, desc, no, desc)
}

func (filterinfo) Close() error {
	return insertFilterinfoSQL.Close()
}
