package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sitemap/link"
	"strings"
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlSet struct {
	Xmlns string `xml:"xmlns,attr"`
	URLS  []loc  `xml:"url"`
}

func main() {
	urlFlag := flag.String("url", "https://gophercises.com", "url that you want to build a sitemap for")
	maxDepth := flag.Int("depth", 10, "the maximum depth of links")
	flag.Parse()

	pages := bfs(*urlFlag, *maxDepth)
	toXML := urlSet{
		Xmlns: xmlns,
	}
	//pages := get(*urlFlag)
	_ = pages
	for _, page := range pages {
		toXML.URLS = append(toXML.URLS, loc{page})
	}
	fmt.Print(xml.Header)
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "  ")
	if err := enc.Encode(toXML); err != nil {
		panic(err)
	}
	/*	we only want internal paths at base-url
		/some-path
		base-url/some-path
		#fragment
		mailto:
	*/
}

func bfs(urlStr string, maxDepth int) []string {
	// key is hashed in map - cool!
	// empty struct uses less memory than other things - fun fact!
	var ret []string
	seen := make(map[string]struct{})
	var queue map[string]struct{}
	enqueue := map[string]struct{}{
		urlStr: {},
	}
	for i := 0; i <= maxDepth; i++ {
		queue, enqueue = enqueue, make(map[string]struct{})
		for nqurl, _ := range queue {
			if _, ok := seen[nqurl]; ok {
				continue
			}
			seen[nqurl] = struct{}{}
			for _, seenLnk := range get(nqurl) {
				enqueue[seenLnk] = struct{}{}
			}

		}
	}

	for seenUrl, _ := range seen {
		ret = append(ret, seenUrl)
	}
	return ret
}

func get(urlStr string) []string {
	resp, err := http.Get(urlStr)
	if err != nil {
		panic(err)
	}

	//MUST CLOSE RESPONSE BODY
	// defer -> whenever this function closes, run this
	// Benefit of using defer: can put code at any point in the code
	// Easy to forget about without defer -> someone else might do a dumb
	// You may skip it on accident / ex: if condition -> return OOPS!
	defer resp.Body.Close()

	reqUrl := resp.Request.URL
	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	base := baseUrl.String()
	fmt.Println(base)
	return filter(hrefs(resp.Body, base), withPrefix(base))
}

func filter(links []string, keepfn func(string) bool) []string {
	var ret []string
	for _, lnk := range links {
		if keepfn(lnk) {
			ret = append(ret, lnk)
		}
	}
	return ret
}

func withPrefix(pfx string) func(string) bool {
	return func(lnk string) bool {
		return strings.HasPrefix(lnk, pfx)
	}
}

func hrefs(body io.Reader, base string) []string {
	links, _ := link.Parse(body)

	var ret []string

	for _, lnk := range links {
		switch {
		case strings.HasPrefix(lnk.Href, "/"):
			ret = append(ret, base+lnk.Href)
		case strings.HasPrefix(lnk.Href, "http"):
			ret = append(ret, lnk.Href)
		}
	}
	return ret
}

/**
1. GET the webpage
2. PARSE the links on the page (use package already made)
3. BUILD proper urls
4. FILTER links that have a different BASE
5. FIND all the pages (BFS)
6. PRINT xml
*/
