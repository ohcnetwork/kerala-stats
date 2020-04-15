package test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/saanuregh/kerala_stats/handler"
	"github.com/saanuregh/kerala_stats/parser"
)

func TestLatest(t *testing.T) {
	e := handler.Build()
	req := httptest.NewRequest(echo.GET, "/latest", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/latest")
	handler.Latest(c)
	var v handler.LatestData
	err := json.Unmarshal(rec.Body.Bytes(), &v)
	if err != nil {
		panic(err)
	}
}

func TestHistory(t *testing.T) {
	e := handler.Build()
	req := httptest.NewRequest(echo.GET, "/history", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/history")
	handler.History(c)
	var v parser.History
	err := json.Unmarshal(rec.Body.Bytes(), &v)
	if err != nil {
		panic(err)
	}
}

func TestTestReport(t *testing.T) {
	e := handler.Build()
	req := httptest.NewRequest(echo.GET, "/test", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/test")
	handler.TestReport(c)
	var v parser.TestReports
	err := json.Unmarshal(rec.Body.Bytes(), &v)
	if err != nil {
		panic(err)
	}
}

func TestLastUpdated(t *testing.T) {
	e := handler.Build()
	req := httptest.NewRequest(echo.GET, "/last", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/last")
	handler.LastUpdated(c)
}
