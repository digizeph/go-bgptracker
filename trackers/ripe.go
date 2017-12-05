package trackers

import (
	"regexp"
)

// GetFirstAndLastTime returns the first and last timestamp of the data dump file from RIPE
func (source RipeSource) GetCollectorInfo() []CollectorInfo {
	// get list of collectors from the front page
	r, _ := regexp.Compile(`http://data\.ris\.ripe.*`)
	collectors := getLinksOnPage(source.RootUrl, r)

	var ch = make(chan CollectorInfo)

	// spin a go routine for each collector
	for _, v := range collectors {
		go crawlCollector(v, regexp.MustCompile(`rrc\d\d`), "", ch)
	}

	// wait for information coming into the channel
	var lst []CollectorInfo
	for range collectors {
		d := <-ch
		lst = append(lst, d)
	}

	return lst
}
