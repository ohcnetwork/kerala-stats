package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

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

func main() {
	lastUpdated := scraper.ScrapeLastUpdated()
	var l LatestHistory
	data, err := ioutil.ReadFile("./latest.json")
	if err != nil {
		log.Panicln(err)
	}
	err = json.Unmarshal(data, &l)
	if err != nil {
		log.Panicln(err)
	}
	if l.LastUpdated != lastUpdated {
		log.Println("Changes found, updating...")
		_histories := scraper.ScrapeHistory()
		histories := Histories{History: _histories, LastUpdated: lastUpdated}
		j, err := json.Marshal(histories)
		if err != nil {
			log.Panicln(err)
		}
		err = ioutil.WriteFile("./histories.json", j, 0644)
		if err != nil {
			log.Panicln(err)
		}
		testReports := TestReports{Reports: scraper.ScrapeTestReport(), LastUpdated: lastUpdated}
		j, err = json.Marshal(testReports)
		if err != nil {
			log.Panicln(err)
		}
		err = ioutil.WriteFile("./testreports.json", j, 0644)
		if err != nil {
			log.Panicln(err)
		}
		latest := len(_histories) - 1
		latestData := LatestHistory{Summary: _histories[latest].Summary, Delta: _histories[latest].Delta, LastUpdated: lastUpdated}
		j, err = json.Marshal(latestData)
		if err != nil {
			log.Panicln(err)
		}
		err = ioutil.WriteFile("./latest.json", j, 0644)
		if err != nil {
			log.Panicln(err)
		}
		s, d := scraper.LatestSummary(&_histories[latest])
		summary := Summary{Summary: s, Delta: d, LastUpdated: lastUpdated}
		j, err = json.Marshal(summary)
		if err != nil {
			log.Panicln(err)
		}
		err = ioutil.WriteFile("./summary.json", j, 0644)
		if err != nil {
			log.Panicln(err)
		}
		log.Println("Updated")
	} else {
		log.Println("No changes")
	}

}
