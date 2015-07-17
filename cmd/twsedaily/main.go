package main

import (
	"flag"
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
var twsecate = flag.String("twsecate", "", "twse cate")

func doCheck(stock *twse.Data) []bool {
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
	flag.Parse()
	recentlyOpened := tradingdays.FindRecentlyOpened(time.Now())
	dailyreportdb := tdb.NewDailyReportDB()
	defer dailyreportdb.Close()

	var stockList = make([]*twse.Data, 0)
	if *twsecate != "" {
		for _, sno := range gettwsecate(*twsecate, recentlyOpened) {
			stockList = append(stockList, twse.NewTWSE(sno, recentlyOpened))
		}
	}

	wg.Add(len(stockList))
	for _, stock := range stockList {
		go func(stock *twse.Data, recentlyOpened time.Time) {
			defer wg.Done()
			runtime.Gosched()
			for i, result := range doCheck(stock) {
				if result {
					if _, err := dailyreportdb.InsertRecode(stock.No, uint64(i), recentlyOpened); err == nil {
						log.Println(filter.AllList[i])
					} else {
						log.Println("InsertRecode Error", stock.No, i, err)
					}
				}
			}
		}(stock, recentlyOpened)
	}
	wg.Wait()
}
