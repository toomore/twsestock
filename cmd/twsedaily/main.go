package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/codegangsta/cli"
	"github.com/toomore/gogrs/cmd/twsereport/filter"
	"github.com/toomore/gogrs/tradingdays"
	"github.com/toomore/gogrs/twse"
	"github.com/toomore/gogrs/utils"
	"github.com/toomore/twsestock/tdb"
)

var (
	flushcache   bool
	otccate      string
	twsecate     string
	updateFilter bool
	wg           sync.WaitGroup
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 4)
}

func doCheck(stock *twse.Data) []bool {
	result := make([]bool, len(filter.AllList))
	for i, filterFunc := range filter.AllList {
		result[i] = filterFunc.CheckFunc(stock)
	}
	return result
}

func gettwsecate(isTwse bool, cate string, date time.Time) []string {
	var l twse.BaseLists
	switch isTwse {
	case true:
		l = twse.NewLists(date)
	default:
		l = twse.NewOTCLists(date)
	}
	var result []string
	for _, s := range l.GetCategoryList(cate) {
		result = append(result, s.No)
	}
	return result
}

func makeStockList(twsecae string, otccate string, recentlyOpened time.Time) []*twse.Data {
	var stockList = make([]*twse.Data, 0)
	if twsecate != "" {
		for _, twsecateno := range strings.Split(twsecate, ",") {
			for _, sno := range gettwsecate(true, twsecateno, recentlyOpened) {
				stockList = append(stockList, twse.NewTWSE(sno, recentlyOpened))
			}
		}
	}

	if otccate != "" {
		for _, otccateno := range strings.Split(otccate, ",") {
			for _, sno := range gettwsecate(false, otccateno, recentlyOpened) {
				stockList = append(stockList, twse.NewOTC(sno, recentlyOpened))
			}
		}
	}
	return stockList
}

func updateFilterInfo() {
	filterinfodb := tdb.NewFilterinfoDB()
	defer filterinfodb.Close()
	for _, v := range filter.AllList {
		fmt.Println(v.No(), v)
		filterinfodb.InsertFilterinfo(v.No(), v.String())
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "twsedaily"
	app.Usage = "Get daily data and into DB."
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "twsecate",
			Value: "",
			Usage: "twsecate no",
		},
		cli.StringFlag{
			Name:  "otccate",
			Value: "",
			Usage: "otccate no",
		},
		cli.BoolFlag{
			Name:  "flushcache",
			Usage: "clear cache",
		},
		cli.BoolFlag{
			Name:  "updatefilter",
			Usage: "update filter",
		},
	}
	app.Action = func(c *cli.Context) {
		if c.NumFlags() == 0 {
			cli.ShowAppHelp(c)
			os.Exit(0)
		}
		flushcache = c.Bool("flushcache")
		otccate = c.String("otccate")
		twsecate = c.String("twsecate")
		updateFilter = c.Bool("updatefilter")
	}
	app.Run(os.Args)

	if flushcache {
		temppath := utils.GetOSRamdiskPath()
		utils.NewHTTPCache(temppath, "utf8").FlushAll()
		log.Println("Clear cache", temppath)
	}

	if updateFilter {
		updateFilterInfo()
	}

	recentlyOpened := tradingdays.FindRecentlyOpened(time.Now())
	dailyreportdb := tdb.NewDailyReportDB()
	defer dailyreportdb.Close()

	stockList := makeStockList(twsecate, otccate, recentlyOpened)

	wg.Add(len(stockList))
	for _, stock := range stockList {
		go func(stock *twse.Data, recentlyOpened time.Time) {
			defer wg.Done()
			runtime.Gosched()
			for i, result := range doCheck(stock) {
				if result {
					if _, err := dailyreportdb.InsertRecode(stock.No, uint64(i), recentlyOpened); err == nil {
						log.Println(stock.No, filter.AllList[i])
					} else {
						log.Println("InsertRecode Error", stock.No, i, err)
					}
				}
			}
		}(stock, recentlyOpened)
	}
	wg.Wait()
}
