package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/saanuregh/kerala_stats/parser"
)

var app *echo.Echo
var historyCache *parser.History
var lastCacheUpdate time.Time

type LatestData struct {
	Kerala      map[string]parser.HistoryDistrict `json:"kerala"`
	LastUpdated string                            `json:"last_updated"`
}

func ParseAndCache() {
	if (historyCache != nil) && (time.Now().Sub(lastCacheUpdate) < (20 * time.Minute)) {
		return
	}
	historyCache = parser.ParseHistory()
	lastCacheUpdate = time.Now()
}

func History(c echo.Context) error {
	ParseAndCache()
	c.Response().Header().Set("Cache-Control", "s-maxage=600, stale-while-revalidate")
	return c.JSON(http.StatusOK, historyCache)
}

func TestReport(c echo.Context) error {
	c.Response().Header().Set("Cache-Control", "s-maxage=600, stale-while-revalidate")
	return c.JSON(http.StatusOK, parser.ParseTestReport())
}

func Latest(c echo.Context) error {
	ParseAndCache()
	c.Response().Header().Set("Cache-Control", "s-maxage=600, stale-while-revalidate")
	return c.JSON(http.StatusOK, &LatestData{
		Kerala:      historyCache.History[len(historyCache.History)-1].Kerala,
		LastUpdated: historyCache.LastUpdated,
	})
}

func LastUpdated(c echo.Context) error {
	c.Response().Header().Set("Cache-Control", "s-maxage=600, stale-while-revalidate")
	return c.String(http.StatusOK, parser.ParseLastUpdated())
}

func Build() (e *echo.Echo) {
	e = echo.New()
	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	e.GET("/history", History)
	e.GET("/latest", Latest)
	e.GET("/test", TestReport)
	e.GET("/last", LastUpdated)
	return
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if app == nil {
		app = Build()
	}
	app.ServeHTTP(w, r)
}
