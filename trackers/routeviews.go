package trackers

import (
	"regexp"
)

// Get the first and last data dump file names from RouteViews collectors
func (source RvSource) GetCollectorInfo() []CollectorInfo {
	// get list of collectors from the front page
	r, _ := regexp.Compile(`/bgpdata`)
	collectors := getLinksOnPage(source.RootUrl, r)

	var ch = make(chan CollectorInfo)

	// spin a go routine for each collector
	for _, v := range collectors {
		go crawlCollector(source.RootUrl+v, regexp.MustCompile(`route-views[a-zA-Z0-9_.-]+/`), "UPDATES/", ch)
	}

	// wait for information coming into the channel
	var lst []CollectorInfo
	for range collectors {
		d := <-ch
		// special case for route-views2 where the path doesn't contain "/bgpdata"
		if d.CollectorName == "" {
			d.CollectorName = "route-views2"
		}
		lst = append(lst, d)
	}

	return lst
}
