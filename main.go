package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sitemap/link"
	"strings"
)

func main() {
	urlFlag := flag.String("url", "https://gophercises.com", "url that you want to build a sitemap for")
	flag.Parse()

	pages := get(*urlFlag)
	_ = pages
	for _, page := range pages {
		fmt.Println(page)
	}

	/*	we only want internal paths at base-url
		/some-path
		base-url/some-path
		#fragment
		mailto:
	*/
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
