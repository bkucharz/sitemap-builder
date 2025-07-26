package crawler

import (
	"fmt"
	"maps"
	"net/http"
	"net/url"
	"sitemap/internal/link"
	"slices"
	"strings"
)

func Fetch(fetchURL *url.URL) ([]*url.URL, error) {
	response, err := http.Get(fetchURL.String())
	if err != nil {
		return nil, fmt.Errorf("failed fetching website: %s", fetchURL.String())
	}
	defer response.Body.Close()

	links, err := link.Parse(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed parsing website content: %s", fetchURL.String())
	}

	return filterURLs(normalizeLinks(links, fetchURL), withSameHost(fetchURL)), nil

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
			urls, err := Fetch(currentURL)
			if err != nil {
				continue
			}
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
