package tdb

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type dailyreport base

var insertRecodeSQL *sql.Stmt

func NewDailyReportDB() *dailyreport {
	table := "dailyreport"
	insertRecodeSQL, _ = conn.Prepare(fmt.Sprintf("Insert into %s(no, filter, timestamp) Values(?,?,?)", table))
	return &dailyreport{
		table: table,
	}
}

func (d dailyreport) InsertRecode(no string, filterno uint64, date time.Time) (sql.Result, error) {
	result, err := insertRecodeSQL.Exec(no, filterno, date)
	log.Println(result.RowsAffected())
	return result, err
}

func (d dailyreport) Close() error {
	return insertRecodeSQL.Close()
}
