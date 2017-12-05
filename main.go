// A crawler that collects BGP updates dump information for first and last
// available data for all collectors in RIPE RRC and RouteViews
package main

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/digizeph/go-bgptracker/trackers"
)

// a global variable that stores the last-crawling timestamp
var lastCrawlTime time.Time

func main() {
	// a map of collector name to collector information
	var collectorsMap = make(map[string]trackers.CollectorInfo)

	// start the web service at the beginning, run as a go routine
	go startServer(collectorsMap, &lastCrawlTime, ":9999")

	// start crawling before the ticker begins
	crawl(collectorsMap)

	// create a ticker that runs for every two minutes
	ticker := time.NewTicker(time.Minute * 2)
	for range ticker.C {
		// when ticker triggers a timeout, crawl the collectors
		crawl(collectorsMap)
	}
}

func crawl(collectorsMap map[string]trackers.CollectorInfo) {
	// update the last-crawling time
	fmt.Println("update update files at time", time.Now())
	lastCrawlTime = time.Now()

	// crawl RouteViews collectors
	fmt.Println("\t RouteViews")
	rv := trackers.RvSource{RootUrl: "http://archive.routeviews.org/"}
	rvCollectors := rv.GetCollectorInfo()
	for _, v := range rvCollectors {
		// update map for web display
		collectorsMap[v.CollectorName] = v
	}

	// crawl RIPE RRC collectors
	fmt.Println("\t RIPE")
	ripe := trackers.RipeSource{RootUrl: "http://data.ris.ripe.net/"}
	ripeCollectors := ripe.GetCollectorInfo()
	for _, v := range ripeCollectors {
		// update map for web display
		collectorsMap[v.CollectorName] = v
	}

	fmt.Println("\t Done")
}

// start http server
func startServer(collectorsMap map[string]trackers.CollectorInfo, timePtr *time.Time, portStr string) {

	h := http.NewServeMux()

	h.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		var collectors []string
		for k := range collectorsMap {
			collectors = append(collectors, k)
		}
		sort.Strings(collectors)

		fmt.Fprintln(w, "<!DOCTYPE html>")
		fmt.Fprintln(w, "<html>")
		fmt.Fprintln(w, "<head>\r<style>\rtable, th, td {\rborder: 1px solid black;\r}\r</style>\r</head>")
		fmt.Fprintln(w, "<body>")

		var tmpl = "<tr>" +
			"<td><a href=\"%s\">%s</a></td>" +
			"<td>%s</td><td>%s</td></tr>\r"
		fmt.Fprintln(w, "<h1>BGP Data Sources First and Last Dump File Time</h1>")
		fmt.Fprintln(w, "<table>")
		fmt.Fprintln(w, "<tr><td>Name</td><td>First</td><td>Last</td></tr>")
		for _, v := range collectors {
			fmt.Fprintf(w, tmpl,
				collectorsMap[v].BaseUrl, collectorsMap[v].CollectorName,
				collectorsMap[v].FirstUrl, collectorsMap[v].LastUrl)
		}
		fmt.Fprintln(w, "</table>")
		fmt.Fprintln(w, "Last updated:", timePtr.String())

		fmt.Fprintln(w, "</body>\r</html>")
	})

	err := http.ListenAndServe(portStr, h)
	log.Fatal(err)
}
