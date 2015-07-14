package tdb

import (
	"database/sql"
	"fmt"
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

func (d dailyreport) InsertRecode(no string, filterno uint64, date time.Time) {
	insertRecodeSQL.Exec(no, filterno, date)
	defer insertRecodeSQL.Close()
}
