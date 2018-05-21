package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// ----------------------------------------------------------------------

// return the URL and query params for the API call
func getURL(page int, t1, t2 time.Time) string {
	url := "https://api.stackexchange.com"
	url += "/2.2/questions?"

	url += fmt.Sprintf("page=%d", page)
	url += fmt.Sprintf("&fromdate=%d", t1.Unix())
	url += fmt.Sprintf("&todate=%d", t2.Unix())

	url += "&pagesize=100"
	url += "&order=desc"
	url += "&sort=creation"
	url += "&tagged=python"
	url += "&site=stackoverflow"

	return url
}

// ----------------------------------------------------------------------

// return the body of the response for the given page and time/date range
func getBody(page int, t1, t2 time.Time) []byte {
	url := getURL(page, t1, t2)

	// Build the request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
	}
	defer resp.Body.Close()

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal("ReadAll: ", readErr)
	}

	return body
}

// ----------------------------------------------------------------------

// we just pull out a couple things from the results.
type SOItem struct {
	CreationDate int64 `json:"creation_date"`
	ViewCount    int   `json:"view_count"`
}

type SOResponse struct {
	Items          []SOItem `json:"Items"`
	HasMore        bool     `json:"has_more"`
	QuotaRemaining int      `json:"quota_remaining"`
}

// ----------------------------------------------------------------------

func daystr(t time.Time) string {
	return fmt.Sprintf("%d %s %d", t.Day(), t.Month().String(), t.Year())
}

// ----------------------------------------------------------------------

func timeFromUnix(unix int64) time.Time {
	return time.Unix(unix, 0).UTC()
}

// ----------------------------------------------------------------------

// Retrieve all the relevant pages, and sum the results for each day
func findbestday(t1 time.Time, t2 time.Time) (string, int64) {
	var results map[string]int64 // use a map to not assume anything about date range.

	results = make(map[string]int64)
	more := true

	for page := 1; more; page++ {
		pagejson := getBody(page, t1, t2)
		var result SOResponse

		json.Unmarshal(pagejson, &result)

		more = result.HasMore

		if result.QuotaRemaining < 100 {
			log.Fatal("PANIC! Our quota is almost spent.")
		}
		fmt.Printf("Page=%d, Quota=%d, has_more=%t, Num Items=%d\n\n", page, result.QuotaRemaining, more, len(result.Items))

		for _, itm := range result.Items {
			t := daystr(timeFromUnix(itm.CreationDate))
			results[t] += int64(itm.ViewCount)
		}
	}

	var maxcount int64 = 0
	var maxday string = "Undefined"
	for key, value := range results {
		fmt.Println(key, " -> ", value)
		if value > maxcount {
			maxday = key
			maxcount = value
		}
	}

	return maxday, maxcount
}

// ----------------------------------------------------------------------

func main() {
	t1, _ := time.Parse(time.RFC822, "01 Mar 18 00:00 UTC")
	//t2, _ := time.Parse(time.RFC822, "01 Mar 18 23:59 UTC") // For testing (7 pages)
	t2, _ := time.Parse(time.RFC822, "31 Mar 18 23:59 UTC") // Note 1522454400 is only start of day of 31st.

	fmt.Println("From ", t1.Format(time.UnixDate), " to ", t2.Format(time.UnixDate))

	day, count := findbestday(t1, t2)

	fmt.Println("The best day was ", day)
	fmt.Println("On that day, the total view count of all questions was ", count)
}
