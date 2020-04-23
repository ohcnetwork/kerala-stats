package scraper

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var districtMap = map[string]string{"1614": "Thiruvananthapuram",
	"1613": "Kollam",
	"1612": "Pathanamthitta",
	"1611": "Alappuzha",
	"1610": "Kottayam",
	"1609": "Idukki",
	"1608": "Ernakulam",
	"1607": "Thrissur",
	"1606": "Palakkad",
	"1605": "Malappuram",
	"1604": "Kozhikode",
	"1603": "Wayanad",
	"1602": "Kannur",
	"1601": "Kasaragod"}

type DistrictInfo struct {
	HospitalObservation int `json:"hospital_obs"`
	HomeObservation     int `json:"home_obs"`
	TotalObservation    int `json:"total_obs"`
	HospitalizedToday   int `json:"hospital_today"`
	Confirmed           int `json:"confirmed"`
	Recovered           int `json:"recovered"`
	Deceased            int `json:"deceased"`
	Active              int `json:"active"`
}

type History struct {
	sync.RWMutex
	Summary map[string]DistrictInfo `json:"summary"`
	Delta   map[string]DistrictInfo `json:"delta"`
	Date    string                  `json:"date"`
}

type TestReport struct {
	Date     string `json:"date"`
	Total    int    `json:"total"`
	Negative int    `json:"negative"`
	Positive int    `json:"positive"`
	Pending  int    `json:"pending"`
	Today    int    `json:"today"`
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func getDoc(source string, referer string, dist ...string) goquery.Document {
	client := &http.Client{}
	var req *http.Request
	if len(dist) > 0 {
		data := url.Values{"district": {dist[0]}, "submit": {"View"}, "vw": {"View"}}
		req, _ = http.NewRequest("POST", source, strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req, _ = http.NewRequest("GET", source, nil)
	}
	req.Host = "dashboard.kerala.gov.in"
	req.Header.Set("Origin", "https://dashboard.kerala.gov.in")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:75.0) Gecko/20100101 Firefox/75.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-GB,en;q=0.5")
	req.Header.Set("Referer", referer)
	req.Header.Set("Connection", "keep-alive")
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		panic(fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status))
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		panic(err)
	}
	return *doc
}

func ScrapeLastUpdated() string {
	url := "https://dashboard.kerala.gov.in/index.php"
	doc := getDoc(url, url)
	s := doc.Find(".breadcrumb-item > i:nth-child(1)").Text()
	s = strings.ToUpper(strings.TrimSpace(strings.Split(s, "Update:")[1]))
	return s
}

func ScrapeTestReport() []TestReport {
	doc := getDoc(
		"https://dashboard.kerala.gov.in/testing-view-public.php",
		"https://dashboard.kerala.gov.in/quar_dst_wise_public.php",
	)
	var row []string
	var rows [][]string
	doc.Find("table > tbody").Each(func(index int, tablehtml *goquery.Selection) {
		tablehtml.Find("tr").Each(func(indextr int, rowhtml *goquery.Selection) {
			rowhtml.Find("td").Each(func(indexth int, tablecell *goquery.Selection) {
				row = append(row, tablecell.Text())
			})
			rows = append(rows, row)
			row = nil
		})
	})
	b := make([]TestReport, len(rows))
	for i := len(rows) - 1; i > -1; i-- {
		b[len(rows)-i-1] = TestReport{
			Date:     rows[i][0],
			Total:    atoi(rows[i][1]),
			Negative: atoi(rows[i][2]),
			Positive: atoi(rows[i][3]),
			Pending:  atoi(rows[i][4]),
			Today:    atoi(rows[i][5]),
		}
	}
	return b
}

func scrapeHistorySingle(b []History, k string, l int, wg *sync.WaitGroup) {
	defer wg.Done()
	url1 := "https://dashboard.kerala.gov.in/dailyreporting-view-public-districtwise.php"
	url2 := "https://dashboard.kerala.gov.in/quar_dst_wise_public.php"
	doc := getDoc(url1, url1, k)
	var row []string
	data1 := make(map[string][]string)
	doc.Find(".table-hover > tbody:nth-child(3)").Each(func(index int, tablehtml *goquery.Selection) {
		tablehtml.Find("tr").Each(func(indextr int, rowhtml *goquery.Selection) {
			rowhtml.Find("td").Each(func(indexth int, tablecell *goquery.Selection) {
				row = append(row, tablecell.Text())
			})
			data1[row[0]] = row[1:]
			row = nil
		})
	})
	doc = getDoc(url2, url2, k)
	data2 := make(map[string][]string)
	doc.Find("table.table:nth-child(2) > tbody:nth-child(3)").Each(func(index int, tablehtml *goquery.Selection) {
		tablehtml.Find("tr").Each(func(indextr int, rowhtml *goquery.Selection) {
			rowhtml.Find("td").Each(func(indexth int, tablecell *goquery.Selection) {
				row = append(row, tablecell.Text())
			})
			if len(row) > 2 {
				data2[row[0]] = row[1:]
			}
			row = nil
		})
	})
	var j, m = 0, 0
	pr1 := []string{"0", "0", "0", "0"}
	pr2 := []string{"0", "0", "0", "0"}
	var pos, dis, act, det, tot, hos, home, tod, dpos, ddis, dact, ddet, dtot, dhos, dhome, dtod = 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0
	for i := 0; i < l; i++ {
		if r1, ok := data1[b[i].Date]; ok {
			pos += atoi(r1[0])
			dis += atoi(r1[1])
			act = atoi(r1[2])
			det += atoi(r1[3])
			dpos = atoi(r1[0])
			ddis = atoi(r1[1])
			ddet = atoi(r1[3])
			if i != 0 {
				dact = act - atoi(pr1[2])
			} else {
				dact = act
			}
			j++
			if j < len(data1) {
				pr1 = r1
			}
		}
		if r2, ok := data2[b[i].Date]; ok {
			tot = atoi(r2[0])
			hos = atoi(r2[1])
			home = atoi(r2[2])
			tod = atoi(r2[3])
			if i != 0 {
				dtot = tot - atoi(pr2[0])
				dhos = hos - atoi(pr2[1])
				dhome = home - atoi(pr2[2])
				dtod = tod - atoi(pr2[3])
			} else {
				dtot = tot
				dhos = hos
				dhome = home
				dtod = tod
			}
			m++
			if m < len(data2) {
				pr2 = r2
			}
		}
		b[i].RLock()
		b[i].Summary[districtMap[k]] = DistrictInfo{
			Confirmed:           pos,
			Recovered:           dis,
			Active:              act,
			Deceased:            det,
			HospitalObservation: hos,
			HomeObservation:     home,
			TotalObservation:    tot,
			HospitalizedToday:   tod,
		}
		b[i].Delta[districtMap[k]] = DistrictInfo{
			Confirmed:           dpos,
			Recovered:           ddis,
			Active:              dact,
			Deceased:            ddet,
			HospitalObservation: dhos,
			HomeObservation:     dhome,
			TotalObservation:    dtot,
			HospitalizedToday:   dtod,
		}
		b[i].RUnlock()
	}
}

func initHistory() ([]History, int) {
	lastUpdated := ScrapeLastUpdated()
	start := "30-01-2020"
	last := strings.Split(lastUpdated, " ")[0]
	t, _ := time.Parse("02-01-2006", start)
	var list []string
	for {
		l := t.Format("02-01-2006")
		if l == last {
			list = append(list, l)
			break
		}
		list = append(list, l)
		t = t.Add(time.Hour * 24)
	}
	b := make([]History, len(list))
	for i, d := range list {
		b[i] = History{Summary: make(map[string]DistrictInfo), Delta: make(map[string]DistrictInfo), Date: d}
		for _, e := range districtMap {
			b[i].Summary[e] = DistrictInfo{}
			b[i].Delta[e] = DistrictInfo{}
		}
	}
	return b, len(list)
}

func ScrapeHistory() []History {
	var wg sync.WaitGroup
	b, n := initHistory()
	for i := 1601; i <= 1614; i++ {
		wg.Add(1)
		go scrapeHistorySingle(b, strconv.Itoa(i), n, &wg)
	}
	wg.Wait()
	return b
}

func LatestSummary(h *History) (DistrictInfo, DistrictInfo) {
	var pos, dis, act, det, tot, hos, home, tod, dpos, ddis, dact, ddet, dtot, dhos, dhome, dtod = 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0
	for _, info := range h.Summary {
		pos += info.Confirmed
		dis += info.Recovered
		act += info.Active
		det += info.Deceased
		tot += info.TotalObservation
		hos += info.HospitalObservation
		home += info.HomeObservation
		tod += info.HospitalizedToday
	}
	for _, info := range h.Delta {
		dpos += info.Confirmed
		ddis += info.Recovered
		dact += info.Active
		ddet += info.Deceased
		dtot += info.TotalObservation
		dhos += info.HospitalObservation
		dhome += info.HomeObservation
		dtod += info.HospitalizedToday
	}
	summary := DistrictInfo{
		Confirmed:           pos,
		Recovered:           dis,
		Active:              act,
		Deceased:            det,
		HospitalObservation: hos,
		HomeObservation:     home,
		TotalObservation:    tot,
		HospitalizedToday:   tod,
	}
	delta := DistrictInfo{
		Confirmed:           dpos,
		Recovered:           ddis,
		Active:              dact,
		Deceased:            ddet,
		HospitalObservation: dhos,
		HomeObservation:     dhome,
		TotalObservation:    dtot,
		HospitalizedToday:   dtod,
	}
	return summary, delta
}
