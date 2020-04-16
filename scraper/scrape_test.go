package scraper

import (
	"testing"
)

func TestScrapeHistory(t *testing.T) {
	ScrapeHistory()
}

func TestScrapeTestReport(t *testing.T) {
	ScrapeTestReport()
}

func TestScrapeLastUpdated(t *testing.T) {
	ScrapeLastUpdated()
}

func TestTotal(t *testing.T) {
	h := ScrapeHistory()
	LatestSummary(&h[len(h)-1])
}

func BenchmarkScrapeHistory(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ScrapeHistory()
	}
}
