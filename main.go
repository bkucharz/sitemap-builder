package main

import (
	"flag"
	"fmt"
	"maps"
	"net/http"
	"net/url"
	"os"
	"sitemap/link"
	"slices"
	"strings"
)

func main() {
	urlFlag := flag.String("url", "https://go.dev/", "the url for a sitemap build")
	depth := flag.Int("depth", 1, "depth of fetching")
	flag.Parse()
	*urlFlag = strings.TrimSuffix(*urlFlag, "/")
	startUrl, err := url.Parse(*urlFlag)
	if err != nil {
		fmt.Printf("Cannot parse url: %v\n", *urlFlag)
		os.Exit(1)
	}
	paths, _ := Crawl(*startUrl, *depth)
	for _, p := range paths {
		fmt.Println(p)
	}

}

func Fetch(fetchUrl *url.URL) []*url.URL {
	response, err := http.Get(fetchUrl.String())
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	links, err := link.Parse(response.Body)
	if err != nil {
		panic(err)
	}

	return filterUrls(normalizeLinks(links, fetchUrl), withSameHost(fetchUrl))

}

func Crawl(startURL url.URL, depth int) ([]string, error) {
	seen := make(map[string]struct{})
	var queue = []*url.URL{&startURL}

	for n := 1; n > 0 && depth > 0; n = len(queue) {
		for i := 0; i < n; i++ {
			currentURL := queue[0]
			queue = queue[1:]
			if _, ok := seen[currentURL.Path]; ok {
				continue
			}
			seen[currentURL.Path] = struct{}{}
			urls := Fetch(currentURL)
			queue = append(queue, urls...)
		}
		depth--
	}
	return slices.Collect(maps.Keys(seen)), nil
}

func normalizeLinks(links []link.Link, base *url.URL) []*url.URL {
	urls := make([]*url.URL, 0, len(links))
	for _, link := range links {
		destUrl, err := url.Parse(link.Href)
		destUrl.Fragment = ""
		destUrl.RawQuery = ""

		if err != nil {
			fmt.Printf("Cannot parse URL: %v\n", link.Href)
			continue
		}

		if destUrl.Host == "" {
			destUrl.Host = base.Host
		}
		if destUrl.Scheme == "" {
			destUrl.Scheme = base.Scheme
		}
		destUrl.Path = strings.TrimSuffix(destUrl.Path, "/")
		urls = append(urls, destUrl)
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
