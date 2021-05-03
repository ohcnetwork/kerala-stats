import { cheerio } from "https://deno.land/x/cheerio@1.0.4/mod.ts";
import {
  DAILYREPORTING_PAGE,
  DATE_REGEX,
  District,
  HOTSPOTS_PAGE,
  INDEX_PAGE,
  QUARANTINED_PAGE,
  TESTING_PAGE,
  VACCINATION_PAGE,
} from "./constants.ts";
import { ERRORS } from "./errors.ts";
import { DistrictFuse, getHTMLString, LSGFuse } from "./utils.ts";
import {
  DistrictInfo,
  History,
  HistoryInfo,
  HotspotsHistory,
  TestReport,
  VacDistrictInfo,
  VacHistory,
  VacHistoryInfo,
} from "./types.d.ts";

export const scrapeLastUpdated = async () => {
  try {
    const html = await getHTMLString(INDEX_PAGE);
    const s = cheerio.load(html)(".breadcrumb-item").text().toUpperCase().trim()
      .split(": ")[1];
    if (s) return s;
    throw ERRORS.SCRAPE_LAST_UPDATED;
  } catch (error) {
    console.error(error);
    throw ERRORS.SCRAPE_LAST_UPDATED;
  }
};

const scrapeTable = (html: string, selector: string, key = true) => {
  const doc = cheerio.load(html);
  let row: string[] = [];
  const res: Record<District, string[]> = {} as Record<District, string[]>;
  doc(selector).find("tr").each(
    (_, rowhtml) => {
      cheerio(rowhtml).find("td").each(
        (idx, tablecell) => {
          row[idx] = cheerio(tablecell).text().trim();
        },
      );
      const k = key
        ? District[row[0] as keyof typeof District]
        : row[0] as District;
      res[k] = row.slice(1);
      row = [];
    },
  );
  return res;
};

export const scrapeTodaysHistory = async (today: string, last?: History) => {
  try {
    const t0 = performance.now();
    const s1 =
      "section.col-lg-6:nth-child(5) > div:nth-child(1) > div:nth-child(2) > div:nth-child(1) > table:nth-child(1) > tbody:nth-child(3)";
    const s2 = "table.table:nth-child(1) > tbody:nth-child(3)";

    const data1 = scrapeTable(await getHTMLString(DAILYREPORTING_PAGE), s1);
    const data2 = scrapeTable(await getHTMLString(QUARANTINED_PAGE), s2);

    const h: History = {
      date: today,
      delta: {} as HistoryInfo,
      summary: {} as HistoryInfo,
    };
    for (const k in District) {
      const d = District[k as keyof typeof District];
      const d1 = data1[d];
      const d2 = data2[d];
      h.summary[d] = {
        confirmed: parseInt(d1[0]),
        recovered: parseInt(d1[1]),
        active: parseInt(d1[2]),
        deceased: parseInt(d1[3]),
        total_obs: parseInt(d2[0]),
        hospital_obs: parseInt(d2[1]),
        home_obs: parseInt(d2[2]),
        hospital_today: parseInt(d2[3]),
      };
      h.delta[d] = Object.keys(h.summary[d])
        .reduce((a, k) => {
          const key = k as keyof DistrictInfo;
          a[key] = h.summary[d][key] - (last?.summary[d][key] ?? 0);
          return a;
        }, {} as DistrictInfo);
    }
    console.log(
      `scrapeTodaysHistory completed in ${performance.now() - t0} ms`,
    );
    return h;
  } catch (error) {
    console.error(error);
    throw ERRORS.SCRAPE_TODAYS_HISTORY;
  }
};

export const scrapeTodaysTestReport = async (today: string) => {
  try {
    const t0 = performance.now();
    const row: string[] = [];
    const html = await getHTMLString(TESTING_PAGE);
    const table = cheerio.load(html)("table > tbody");
    for (const rowhtml of table.find("tr").toArray()) {
      const el = cheerio(rowhtml);
      const date = el.first().text();
      if (date.match(DATE_REGEX)?.toString() == today) {
        el.find("td").each((_, tablecell) => {
          row.push(cheerio(tablecell).text());
        });
        const b: TestReport = {
          date: row[0],
          total: parseInt(row[1]),
          today: parseInt(row[2]),
          positive: parseInt(row[4]),
          today_positive: parseInt(row[5]),
        };
        console.log(
          `scrapeTodaysTestReport completed in ${performance.now() - t0} ms`,
        );
        return b;
      }
    }
    throw ERRORS.SCRAPE_TODAYS_TEST_REPORT_DATE;
  } catch (error) {
    console.error(error);
    throw ERRORS.SCRAPE_TODAYS_TEST_REPORT;
  }
};

export const scrapeHotspotsHistory = async (today: string) => {
  try {
    const t0 = performance.now();
    const html = await getHTMLString(HOTSPOTS_PAGE);
    const b: HotspotsHistory = { hotspots: [], date: today };
    let row: string[] = [];
    cheerio.load(html)("div.card:nth-child(1) > div:nth-child(2) > table:nth-child(1) > tbody:nth-child(2)").each(
      (_, tablehtml) => {
        cheerio(tablehtml).find("tr").each((_, rowhtml) => {
          cheerio(rowhtml).find("td").each((_, tablecell) => {
            row.push(cheerio(tablecell).text().trim());
          });
          if (row.length !== 0 && row[1] && row[2]) {
            row[2] = row[2].replace("Muncipality", "(M)");
            row[2] = row[2].replace("©", "(C)");
            if (row[2] === "Koothuparamba (M)") {
              row[2] = "Kuthuparambu (M)";
            }
            if (row[2] === "Mattanur (M)") {
              row[2] = "Mattannoor (M)";
            }
            if (row[2] === "Maloor") {
              row[2] = "Malur";
            }
            if (row[2] === "Changanacherry (M)") {
              row[2] = "Changanassery (M)";
            }
            if (row[2] === "District Hospital") {
              row[2] = "Marutharoad";
            }
            if (row[2] === "Neduveli") {
              row[2] = "Vembayam";
            }

            const d = DistrictFuse.search(row[1], { limit: 1 });
            const s = LSGFuse[d[0].item as District].search(row[2], {
              limit: 1,
            });
            // if s.Score < 60 || d.Score < 60 {
            // 	log.Printf("found innaccurrate matching for %v:%v %v:%v\n", row[1], d.Match, row[2], s.Match)
            // }
            b.hotspots.push({
              district: d[0].item,
              lsgd: s[0].item,
              wards: row[3].replace(/\n\s*/gi, " ").replace(/\s{2,}/gi, " "),
            });
          }
          row = [];
        });
      },
    );
    if (b.hotspots.length < 1) {
      throw ERRORS.SCRAPE_HOTSPOTS_HISTORY;
    }
    console.log(
      `scrapeHotspotsHistory completed in ${performance.now() - t0} ms`,
    );
    return b;
  } catch (error) {
    console.error(error);
    throw ERRORS.SCRAPE_HOTSPOTS_HISTORY;
  }
};

export const scrapeVaccinationHistory = async (
  today: string,
  last?: VacHistory,
) => {
  try {
    const t0 = performance.now();
    const data = scrapeTable(
      await getHTMLString(VACCINATION_PAGE),
      "table.table-hover:nth-child(1) > tbody:nth-child(2)",
      false,
    );
    const h: VacHistory = {
      date: today,
      delta: {} as VacHistoryInfo,
      summary: {} as VacHistoryInfo,
    };
    for (const k in District) {
      const d = District[k as keyof typeof District];
      const x = data[d];
      h.summary[d] = {
        cvc_public: parseInt(x[0]),
        cvc_private: parseInt(x[1]),
        cvc_total: parseInt(x[2]),
        hcw_dose1: parseInt(x[3]),
        hcw_dose2: parseInt(x[4]),
        flw_other_dose1: parseInt(x[5]),
        flw_other_dose2: parseInt(x[6]),
        flw_polling_dose1: parseInt(x[7]),
        flw_polling_dose2: parseInt(x[8]),
        age_appropriate_dose1: parseInt(x[9]),
        age_appropriate_dose2: parseInt(x[10]),
        tot_vaccinations: parseInt(x[11]),
        tot_person_vaccinations: parseInt(x[12]),
        second_dose: parseInt(x[13]),
      };
      h.delta[d] = Object.keys(h.summary[d])
        .reduce((a, k) => {
          const key = k as keyof VacDistrictInfo;
          a[key] = h.summary[d][key] - (last?.summary[d][key] ?? 0);
          return a;
        }, {} as VacDistrictInfo);
    }
    console.log(
      `scrapeVaccinationHistory completed in ${performance.now() - t0} ms`,
    );
    return h;
  } catch (error) {
    console.error(error);
    throw ERRORS.SCRAPE_VACCINATION_HISTORY;
  }
};
