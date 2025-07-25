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
	urlFlag := flag.String("url", "https://go.dev/", "the url for a sitemap build")
	flag.Parse()

	response, err := http.Get(*urlFlag)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	links, err := link.Parse(response.Body)
	if err != nil {
		panic(err)
	}

	reqUrl := response.Request.URL
	startUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}

	urls := filterUrls(getUrls(links, startUrl), withSameHost(startUrl))
	for _, url := range urls {
		fmt.Printf("%#v\n", url.String())
	}
}

func getUrls(links []link.Link, baseUrl *url.URL) []*url.URL {
	var urls []*url.URL
	for _, link := range links {
		url, err := url.Parse(link.Href)
		url.Fragment = ""
		url.RawQuery = ""

		if err != nil {
			fmt.Printf("Cannot parse URL: %v\n", link.Href)
			continue
		}

		if url.Host == "" {
			url.Host = baseUrl.Host
		}
		if url.Scheme == "" {
			url.Scheme = baseUrl.Scheme
		}
		urls = append(urls, url)
	}
	return urls
}

func filterUrls(urls []*url.URL, keepFn func(*url.URL) bool) []*url.URL {
	var filtered []*url.URL
	for _, url := range urls {
		if keepFn(url) {
			filtered = append(filtered, url)
		}
	}
	return filtered
}

func withSameHost(baseUrl *url.URL) func(*url.URL) bool {
	return func(url *url.URL) bool {
		return (url.Host == baseUrl.Host && strings.HasPrefix(url.Scheme, "http"))
	}
}
