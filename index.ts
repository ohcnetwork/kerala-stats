import {
  HISTORIES_FILE,
  HOTSPOT_HISTORIES_FILE,
  HOTSPOT_LATEST_FILE,
  LATEST_FILE,
  SUMMARY_FILE,
  TEST_REPORTS_FILE,
  VAC_HISTORIES_FILE,
  VAC_LATEST_FILE,
  VAC_SUMMARY_FILE,
} from "./src/constants.ts";
import {
  handler,
  lastUpdatedToDate,
  readJson,
  summarizeDistrictInfo,
  writeJson,
} from "./src/utils.ts";
import {
  scrapeHotspotsHistory,
  scrapeLastUpdated,
  scrapeTodaysHistory,
  scrapeTodaysTestReport,
  scrapeVaccinationHistory,
} from "./src/scrapers.ts";
import {
  Histories,
  History,
  HotspotsHistories,
  LatestHistory,
  LatestHotspotsHistory,
  Summary,
  TestReports,
  VacHistories,
  VacHistory,
  VacLatestHistory,
  VacSummary,
} from "./src/types.d.ts";
import { ERRORS } from "./src/errors.ts";

const handleHistories = async (lastUpdated: string) => {
  try {
    const histories = await readJson(
      HISTORIES_FILE,
    ) as Histories;
    const last = histories.histories.length - 1;
    const date = lastUpdatedToDate(lastUpdated);

    let b: History;
    if (date === histories.histories[last]?.date) {
      b = await scrapeTodaysHistory(date, histories.histories[last - 1]);
      histories.histories[last] = b;
      console.log("history replaced");
    } else {
      b = await scrapeTodaysHistory(date, histories.histories[last]);
      histories.histories.push(b);
      console.log("history appended");
    }
    histories.last_updated = lastUpdated;
    await writeJson(HISTORIES_FILE, histories);
    console.log("histories written");

    const latestData: LatestHistory = {
      summary: b.summary,
      delta: b.delta,
      last_updated: lastUpdated,
    };
    await writeJson(LATEST_FILE, latestData);
    console.log("latest written");

    const summaryData: Summary = {
      summary: summarizeDistrictInfo(Object.values(b.summary)),
      delta: summarizeDistrictInfo(Object.values(b.delta)),
      last_updated: lastUpdated,
    };
    await writeJson(SUMMARY_FILE, summaryData);
    console.log("summary written");
  } catch (error) {
    console.error(error);
    throw ERRORS.HANDLE_HISTORIES;
  }
};

const handleTestReports = async (lastUpdated: string) => {
  try {
    const testReports = await readJson(
      TEST_REPORTS_FILE,
    ) as TestReports;
    const last = testReports.reports.length - 1;
    const date = lastUpdatedToDate(lastUpdated);

    const latest = await scrapeTodaysTestReport(date);
    if (date === testReports.reports[last]?.date) {
      testReports.reports[last] = latest;
      console.log("test report replaced");
    } else {
      testReports.reports.push(latest);
      console.log("test report appended");
    }
    testReports.last_updated = lastUpdated;
    await writeJson(TEST_REPORTS_FILE, testReports);
    console.log("test reports written");
  } catch (error) {
    console.error(error);
    throw ERRORS.HANDLE_TEST_REPORTS;
  }
};

const handleHotspotsHistories = async (lastUpdated: string) => {
  try {
    const hhistories = await readJson(
      HOTSPOT_HISTORIES_FILE,
    ) as HotspotsHistories;
    const last = hhistories.histories.length - 1;
    const date = lastUpdatedToDate(lastUpdated);

    const hh = await scrapeHotspotsHistory(date);
    if (date === hhistories.histories[last]?.date) {
      hhistories.histories[last] = hh;
      console.log("hotspot history replaced");
    } else {
      hhistories.histories.push(hh);
      console.log("hotspot history appended");
    }
    hhistories.last_updated = lastUpdated;
    await writeJson(HOTSPOT_HISTORIES_FILE, hhistories);
    console.log("hotspots histories written");

    const latestHotspotData: LatestHotspotsHistory = {
      hotspots: hh.hotspots,
      last_updated: lastUpdated,
    };
    await writeJson(HOTSPOT_LATEST_FILE, latestHotspotData);
    console.log("hotspots latest written");
  } catch (error) {
    console.error(error);
    throw ERRORS.HANDLE_HOTSPOT_HISTORIES;
  }
};

const handleVaccinationHistories = async (lastUpdated: string) => {
  try {
    const histories = await readJson(
      VAC_HISTORIES_FILE,
    ) as VacHistories;
    const last = histories.histories.length - 1;
    const date = lastUpdatedToDate(lastUpdated);

    let b: VacHistory;
    if (date === histories.histories[last]?.date) {
      b = await scrapeVaccinationHistory(date, histories.histories[last - 1]);
      histories.histories[last] = b;
      console.log("vaccination history replaced");
    } else {
      b = await scrapeVaccinationHistory(date, histories.histories[last]);
      histories.histories.push(b);
      console.log("vaccination history appended");
    }
    histories.last_updated = lastUpdated;
    await writeJson(VAC_HISTORIES_FILE, histories);
    console.log("vaccination histories written");

    const latestData: VacLatestHistory = {
      summary: b.summary,
      delta: b.delta,
      last_updated: lastUpdated,
    };
    await writeJson(VAC_LATEST_FILE, latestData);
    console.log("vaccination latest written");

    const summaryData: VacSummary = {
      summary: summarizeDistrictInfo(Object.values(b.summary)),
      delta: summarizeDistrictInfo(Object.values(b.delta)),
      last_updated: lastUpdated,
    };
    await writeJson(VAC_SUMMARY_FILE, summaryData);
    console.log("vaccination summary written");
  } catch (error) {
    console.error(error);
    throw ERRORS.HANDLE_VACCINATION_HISTORIES;
  }
};

const main = async () => {
  const t0 = performance.now();
  const lastUpdated = await scrapeLastUpdated();
  console.log(`last updated on ${lastUpdated}`);
  await handler(handleHistories, lastUpdated);
  await handler(handleTestReports, lastUpdated);
  await handler(handleHotspotsHistories, lastUpdated);
  await handler(handleVaccinationHistories, lastUpdated);
  console.log(
    `scraping completed in ${performance.now() - t0} ms`,
  );
};

await handler(main);
