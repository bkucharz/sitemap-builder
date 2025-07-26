package main

import (
	"encoding/xml"
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

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlSet struct {
	URLs  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

func main() {
	urlFlag := flag.String("url", "https://go.dev/", "the url for a sitemap build")
	depth := flag.Int("depth", 1, "depth of fetching")
	flag.Parse()
	*urlFlag = strings.TrimSuffix(*urlFlag, "/")

	startURL, err := url.Parse(*urlFlag)
	if err != nil {
		fmt.Printf("Cannot parse url: %v\n", *urlFlag)
		os.Exit(1)
	}

	paths, _ := Crawl(*startURL, *depth)
	toXML := urlSet{
		Xmlns: xmlns,
	}
	for _, p := range paths {
		toXML.URLs = append(toXML.URLs, loc{startURL.Host + p})
	}

	out, err := xml.MarshalIndent(toXML, "", " ")
	if err != nil {
		fmt.Printf("Error while creating XML\n")
		os.Exit(1)
	}
	fmt.Println(xml.Header + string(out))
}

func Fetch(fetchURL *url.URL) []*url.URL {
	response, err := http.Get(fetchURL.String())
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	links, err := link.Parse(response.Body)
	if err != nil {
		panic(err)
	}

	return filterURLs(normalizeLinks(links, fetchURL), withSameHost(fetchURL))

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
		destURL, err := url.Parse(link.Href)
		destURL.Fragment = ""
		destURL.RawQuery = ""

		if err != nil {
			fmt.Printf("Cannot parse URL: %v\n", link.Href)
			continue
		}

		if destURL.Host == "" {
			destURL.Host = base.Host
		}
		if destURL.Scheme == "" {
			destURL.Scheme = base.Scheme
		}
		destURL.Path = strings.TrimSuffix(destURL.Path, "/")
		urls = append(urls, destURL)
	}
	return urls
}

func filterURLs(urls []*url.URL, keepFn func(*url.URL) bool) []*url.URL {
	filtered := make([]*url.URL, 0, len(urls))
	for _, u := range urls {
		if keepFn(u) {
			filtered = append(filtered, u)
		}
	}
	return filtered
}

func withSameHost(baseURL *url.URL) func(*url.URL) bool {
	return func(url *url.URL) bool {
		return (url.Host == baseURL.Host && strings.HasPrefix(url.Scheme, "http"))
	}
}
