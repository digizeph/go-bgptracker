package trackers

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"regexp"
	"time"
)

// Helper function to pull the href attribute from a Token
func GetHref(t html.Token) (ok bool, href string) {
	// Iterate over all of the Token's attributes until we find an "href"
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}

	// "bare" return will return the variables (ok, href) as defined in
	// the function definition
	return
}

func getLinksOnPage(url string, regexp *regexp.Regexp) []string {

	//, _ := http.Get("http://data.ris.ripe.net/")
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("ERROR: Failed to crawl \"" + url + "\"")
		return nil
	}

	b := resp.Body
	defer b.Close() // close Body when the function returns
	z := html.NewTokenizer(b)

	var links []string

outer:
	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			fmt.Errorf("error getting html")
			break outer
		case tt == html.StartTagToken:
			t := z.Token()

			// Check if the token is an <a> tag
			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}

			// Extract the href value, if there is one
			ok, url := GetHref(t)
			if !ok {
				continue
			}

			if regexp == nil {
				links = append(links, url)
			} else if regexp.MatchString(url) {
				links = append(links, url)
			}
		}
	}

	return links
}

type UrlTime struct {
	url string
	t   time.Time
}

type ByTime []UrlTime

func (s ByTime) Len() int {
	return len(s)
}
func (s ByTime) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByTime) Less(i, j int) bool {
	return s[i].t.Before(s[j].t)
}
