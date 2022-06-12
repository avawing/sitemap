package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"sitemap/link"
	"strings"
)

func main() {
	urlFlag := flag.String("url", "https://gophercises.com", "url that you want to build a sitemap for")
	flag.Parse()

	resp, err := http.Get(*urlFlag)
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

	links, _ := link.Parse(resp.Body)

	var hrefs []string

	for _, lnk := range links {
		switch {
		case strings.HasPrefix(lnk.Href, "/"):
			hrefs = append(hrefs, base+lnk.Href)
		case strings.HasPrefix(lnk.Href, "http"):
			hrefs = append(hrefs, lnk.Href)
		}
	}

	for _, href := range hrefs {
		fmt.Println(href)
	}
	/*	we only want internal paths at base-url
		/some-path
		base-url/some-path
		#fragment
		mailto:
	*/
}

/**
1. GET the webpage
2. PARSE the links on the page (use package already made)
3. BUILD proper urls
4. FILTER links that have a different BASE
5. FIND all the pages (BFS)
6. PRINT xml
*/
