package test

import (
	"testing"

	"github.com/saanuregh/kerala_stats/parser"
)

func TestParseHistory(t *testing.T) {
	parser.ParseHistory()
}

func TestParseTestReport(t *testing.T) {
	parser.ParseTestReport()
}

func TestParseLastUpdated(t *testing.T) {
	parser.ParseTestReport()
}

func BenchmarkParseHistory(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parser.ParseHistory()
	}
}
