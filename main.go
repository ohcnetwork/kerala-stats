package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"github.com/coronasafe/kerala_stats/scraper"
)

type Histories struct {
	History     []scraper.History `json:"histories"`
	LastUpdated string            `json:"last_updated"`
}

type TestReports struct {
	Reports     []scraper.TestReport `json:"reports"`
	LastUpdated string               `json:"last_updated"`
}

type LatestHistory struct {
	Summary     map[string]scraper.DistrictInfo `json:"summary"`
	Delta       map[string]scraper.DistrictInfo `json:"delta"`
	LastUpdated string                          `json:"last_updated"`
}

type Summary struct {
	Summary     scraper.DistrictInfo `json:"summary"`
	Delta       scraper.DistrictInfo `json:"delta"`
	LastUpdated string               `json:"last_updated"`
}

func writeJSON(v interface{}, filename string) {
	j, err := json.Marshal(v)
	if err != nil {
		log.Panicln(err)
	}
	err = ioutil.WriteFile(filename, j, 0644)
	if err != nil {
		log.Panicln(err)
	}
}

func main() {
	log.Println("started")
	start := time.Now()
	lastUpdated := scraper.ScrapeLastUpdated()
	_histories := scraper.ScrapeHistory()
	histories := Histories{History: _histories, LastUpdated: lastUpdated}
	writeJSON(histories, "./histories.json")
	testReports := TestReports{Reports: scraper.ScrapeTestReport(), LastUpdated: lastUpdated}
	writeJSON(testReports, "./testreports.json")
	latest := len(_histories) - 1
	latestData := LatestHistory{Summary: _histories[latest].Summary, Delta: _histories[latest].Delta, LastUpdated: lastUpdated}
	writeJSON(latestData, "./latest.json")
	s, d := scraper.LatestSummary(&_histories[latest])
	summary := Summary{Summary: s, Delta: d, LastUpdated: lastUpdated}
	writeJSON(summary, "./summary.json")
	log.Printf("completed in %v", time.Now().Sub(start))
}
