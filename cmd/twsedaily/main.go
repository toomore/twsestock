package main

import (
	"log"
	"time"

	"github.com/toomore/gogrs/cmd/twsereport/filter"
	"github.com/toomore/gogrs/tradingdays"
	"github.com/toomore/gogrs/twse"
)

func doCheck(no string, recentlyOpened time.Time) []bool {
	stock := twse.NewTWSE(no, recentlyOpened)
	result := make([]bool, len(filter.AllList))
	for i, filterFunc := range filter.AllList {
		result[i] = filterFunc.CheckFunc(stock)
	}
	return result
}

func main() {
	recentlyOpened := tradingdays.FindRecentlyOpened(time.Now())
	for i, result := range doCheck("1453", recentlyOpened) {
		if result {
			log.Println(filter.AllList[i])
		}
	}
}
