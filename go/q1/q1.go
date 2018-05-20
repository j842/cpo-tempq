package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func getURL(page int, t1, t2 time.Time) string {
	url := fmt.Sprintf("https://api.stackexchange.com/2.2/questions?page=%d&pagesize=99&fromdate=%d&todate=%d&tagged=python&site=stackoverflow",
		page, t1.Unix(), t2.Unix())
	return url
}

func getBody(page int, t1, t2 time.Time) []byte {
	url := getURL(page, t1, t2)

	// Build the request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		os.Exit(1)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal("ReadAll: ", readErr)
		os.Exit(1)
	}

	return body
}

type SOItem struct {
	CreationDate int64 `json:"creation_date"`
	ViewCount    int   `json:"view_count"`
}

type SOResponse struct {
	Items          []SOItem `json:"Items"`
	HasMore        bool     `json:"has_more"`
	QuotaRemaining int      `json:"quota_remaining"`
}

func tallyho(t1 time.Time, t2 time.Time) {
	pagejson := getBody(1, t1, t2)
	//var result map[string]interface{}
	var result SOResponse

	json.Unmarshal(pagejson, &result)

	fmt.Printf("Quota=%d, has_more=%t, Num Items=%d\n\n", result.QuotaRemaining, result.HasMore, len(result.Items))

	for i, itm := range result.Items {
		fmt.Printf("%3d Count=%5d, creation=%d\n", i+1, itm.ViewCount, itm.CreationDate)
	}

}

func main() {
	t1, _ := time.Parse(time.RFC822, "01 Mar 18 00:00 UTC")
	t2, _ := time.Parse(time.RFC822, "31 Mar 18 23:59 UTC") // Note 1522454400 is only start of day of 31st.

	fmt.Println("From ", t1.Format(time.UnixDate), " to ", t2.Format(time.UnixDate))
	fmt.Println("From ", t1.Unix(), " to ", t2.Unix(), " (unix time)")

	tallyho(t1, t2)
}
