package main

import (
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/toomore/gogrs/cmd/twsereport/filter"
	"github.com/toomore/gogrs/tradingdays"
	"github.com/toomore/gogrs/twse"
	"github.com/toomore/twsestock/tdb"
)

var wg sync.WaitGroup

func doCheck(stock *twse.Data, recentlyOpened time.Time) []bool {
	result := make([]bool, len(filter.AllList))
	for i, filterFunc := range filter.AllList {
		result[i] = filterFunc.CheckFunc(stock)
	}
	return result
}

func gettwsecate(cate string, date time.Time) []string {
	l := twse.NewLists(date)
	var result []string
	for _, s := range l.GetCategoryList(cate) {
		result = append(result, s.No)
	}
	return result
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 4)
}

func main() {
	recentlyOpened := tradingdays.FindRecentlyOpened(time.Now())
	dailyreportdb := tdb.NewDailyReportDB()
	defer dailyreportdb.Close()

	for _, sno := range gettwsecate("13", recentlyOpened) {
		wg.Add(1)
		go func(sno string, recentlyOpened time.Time) {
			defer wg.Done()
			runtime.Gosched()
			stock := twse.NewTWSE(sno, recentlyOpened)
			for i, result := range doCheck(stock, recentlyOpened) {
				if result {
					dailyreportdb.InsertRecode(sno, uint64(i), recentlyOpened)
					log.Println(filter.AllList[i])
				}
			}
		}(sno, recentlyOpened)
	}
	wg.Wait()
}
