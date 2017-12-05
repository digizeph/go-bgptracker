package trackers

import (
	"regexp"
	"sort"
	"strings"
	"time"
)

// collector's information struct
type CollectorInfo struct {
	CollectorName string // name of the collector
	BaseUrl       string // base url of the collector
	FirstUrl      string // url of the first available data dump
	LastUrl       string // url of the last available data dump
}

// interface for both RIPE and RouteViews collector
type DataSource interface {
	GetCollectorInfo() []CollectorInfo
}

type RipeSource struct {
	RootUrl string
}

type RvSource struct {
	RootUrl string
}

// crawl a month's data, return the earliest url if isFirst==true, otherwise return latest url
func crawlMonth(baseurl string, isFirst bool, ch chan<- string) {
	// get all links to data dump files in the current month
	r, _ := regexp.Compile(`updates\.(\d\d\d\d\d\d\d\d\.\d\d\d\d)`)
	var links = getLinksOnPage(baseurl, r)

	if len(links) > 0 {
		// sort links by time
		var lst []UrlTime
		for _, v := range links {
			t, _ := time.Parse("20060102.1504", r.FindStringSubmatch(v)[1])
			lst = append(lst, UrlTime{v, t})
		}
		sort.Sort(ByTime(lst))
		// return proper link based on isFirst value
		if isFirst {
			ch <- lst[0].url
		} else {
			ch <- lst[len(lst)-1].url
		}
	} else {
		ch <- ""
	}
}

// crawl a collector's information
// inputs:
//  baseurl -> base url
//  nameRegex -> regular expression to identify collector's name (e.g. route-views3)
//  pathSuffix -> string suffix of the collector path before the data
//  ch -> the channel to return data back to
func crawlCollector(baseurl string, nameRegex *regexp.Regexp, pathSuffix string, ch chan<- CollectorInfo) {
	// append "/" to path
	if !strings.HasSuffix(baseurl, "/") {
		baseurl += "/"
	}

	// get links to monthly data in the collector
	r, _ := regexp.Compile(`\d\d\d\d\.\d\d/`)
	var links = getLinksOnPage(baseurl, r)

	// sort links by time
	var lst []UrlTime
	for _, v := range links {
		t, _ := time.Parse("2006.01/", v)
		if t.After(time.Now()) {
			continue
		}
		lst = append(lst, UrlTime{v + pathSuffix, t})
	}
	sort.Sort(ByTime(lst))

	// get the first and last month
	firstMonth := lst[0]
	lastMonth := lst[len(lst)-1]

	// crawl the first and last month of data in separate go routines
	c1 := make(chan string)
	c2 := make(chan string)
	go crawlMonth(baseurl+firstMonth.url, true, c1)
	go crawlMonth(baseurl+lastMonth.url, false, c2)

	// wait for both month's data back
	var first string
	var last string
	for i := 0; i < 2; i++ {
		select {
		case v1 := <-c1:
			first = v1
		case v2 := <-c2:
			last = v2
		}
	}

	// deal with empty first month case
	for first == "" {
		// go to next month
		firstMonth.t = firstMonth.t.AddDate(0, 1, 0)
		firstMonth.url = firstMonth.t.Format("2006.01/" + pathSuffix)
		go crawlMonth(baseurl+firstMonth.url, true, c1)
		first = <-c1
	}

	// deal with empty last month case
	for last == "" {
		// go to next month
		lastMonth.t = lastMonth.t.AddDate(0, -1, 0)
		lastMonth.url = lastMonth.t.Format("2006.01/" + pathSuffix)
		go crawlMonth(baseurl+lastMonth.url, false, c2)
		last = <-c2
	}

	// return the CollectorInfo with associate information
	ch <- CollectorInfo{nameRegex.FindString(baseurl), baseurl, first, last}
}
