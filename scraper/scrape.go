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

func getDoc(source string, dist ...string) goquery.Document {
	var res *http.Response
	var err error
	if len(dist) > 0 {
		res, err = http.PostForm(source, url.Values{"dist": {dist[0]}})
	} else {
		res, err = http.Get(source)
	}
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
	doc := getDoc("https://dashboard.kerala.gov.in/quarantine-view-public.php")
	s := doc.Find(".breadcrumb-item > small:nth-child(1)").Text()
	s = strings.ToUpper(strings.TrimSpace(strings.Split(s, "Update:")[1]))
	return s
}

func ScrapeTestReport() []TestReport {
	doc := getDoc("https://dashboard.kerala.gov.in/testing-view-public.php")
	var row []string
	var rows [][]string
	doc.Find(".table > tbody:nth-child(3)").Each(func(index int, tablehtml *goquery.Selection) {
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
	doc := getDoc("https://dashboard.kerala.gov.in/ajax_disp_dist_list.php", k)
	data1 := strings.Split(strings.TrimSpace(strings.Replace(strings.Replace(doc.Find("body").Text(), " ", "", -1), "						", "", -1)), "\n\n\n\n")
	for i := 0; i < len(data1); i++ {
		data1[i] = strings.Replace(data1[i], "\n", " ", -1)
	}
	data1 = data1[:len(data1)-1]
	doc = getDoc("https://dashboard.kerala.gov.in/ajax_quarantine_dist_list.php", k)
	data2 := strings.Split(strings.TrimSpace(strings.Replace(strings.Replace(doc.Find("body").Text(), " ", "", -1), "						", "", -1)), "\n\n\n")
	for i := 0; i < len(data2); i++ {
		data2[i] = strings.Replace(data2[i], "\n", " ", -1)
	}
	var j, m = 0, 0
	r1 := strings.Split(data1[j], " ")
	r2 := strings.Split(data2[m], " ")
	var pos, dis, act, det, tot, hos, home, tod, dpos, ddis, dact, ddet, dtot, dhos, dhome, dtod = 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0
	for i := 0; i < l; i++ {
		if r1[0] == b[i].Date {
			pos += atoi(r1[1])
			dis += atoi(r1[2])
			act = atoi(r1[3])
			det += atoi(r1[4])
			dpos = atoi(r1[1])
			ddis = atoi(r1[2])
			ddet = atoi(r1[4])
			if i != 0 {
				dact = act - atoi(r1[3])
			}
			j++
			if j < len(data1) {
				r1 = strings.Split(data1[j], " ")
			}
		}
		if r2[0] == b[i].Date {
			tot = atoi(r2[1])
			hos = atoi(r2[2])
			home = atoi(r2[3])
			tod = atoi(r2[4])
			if i != 0 {
				dtot = tot - atoi(r2[1])
				dhos = hos - atoi(r2[2])
				dhome = home - atoi(r2[3])
				dtod = tod - atoi(r2[4])
			}
			m++
			if m < len(data2) {
				r2 = strings.Split(data2[m], " ")
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
