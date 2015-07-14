package tdb

import (
	"database/sql"
	"fmt"
	"time"
)

type dailyreport struct {
	db    *sql.DB
	table string
}

func NewDailyReportDB() *dailyreport {
	return &dailyreport{
		db:    conn,
		table: "dailyreport",
	}
}

func (d dailyreport) InsertRecode(no string, filterno uint64, date time.Time) {
	stmt, _ := d.db.Prepare(fmt.Sprintf("Insert into %s(no, filter, timestamp) Values(?,?,?)", d.table))
	stmt.Exec(no, filterno, date)
}
