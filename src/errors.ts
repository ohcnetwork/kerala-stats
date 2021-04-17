export const ERRORS = {
  "SCRAPE_LAST_UPDATED": new Error("error scraping last updated"),
  "SCRAPE_TODAYS_HISTORY": new Error("error scraping today's history"),
  "SCRAPE_TODAYS_TEST_REPORT": new Error("error scraping today's test report"),
  "SCRAPE_TODAYS_TEST_REPORT_DATE": new Error(
    "no test report matching the date found",
  ),
  "SCRAPE_HOTSPOTS_HISTORY": new Error("error scraping hotspot table"),
  "SCRAPE_VACCINATION_HISTORY": new Error("error scraping vaccination history"),
  "WEBHOOK_MESSAGE": new Error("error sending webhook message"),
  "WEBHOOK_NO_URL": new Error("could not find discord webhook url"),
  "HANDLE_HISTORIES": new Error("error producing histories files"),
  "HANDLE_TEST_REPORTS": new Error("error producing test report files"),
  "HANDLE_HOTSPOT_HISTORIES": new Error(
    "error producing hotspot histories files",
  ),
  "HANDLE_VACCINATION_HISTORIES": new Error(
    "error producing vaccination histories files",
  ),
};
