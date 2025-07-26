package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/bkucharz/sitemap-builder/internal/crawler"
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlSet struct {
	URLs  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

type Config struct {
	URL      *url.URL
	MaxDepth int
}

func main() {
	conf, err := parseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Config error: %v\n", err)
		os.Exit(1)
	}

	paths, _ := crawler.Crawl(*conf.URL, conf.MaxDepth)
	out := generateXML(paths, *conf)
	fmt.Println(out)
}

func generateXML(paths []string, conf Config) string {
	toXML := urlSet{
		Xmlns: xmlns,
	}
	for _, p := range paths {
		toXML.URLs = append(toXML.URLs, loc{conf.URL.String() + p})
	}

	out, err := xml.MarshalIndent(toXML, "", " ")
	if err != nil {
		fmt.Printf("Error while creating XML\n")
		os.Exit(1)
	}
	return string(out)
}
func parseFlags() (*Config, error) {
	urlFlag := flag.String("url", "https://go.dev/", "the url for a sitemap build")
	depth := flag.Int("depth", 1, "depth of fetching")
	flag.Parse()

	startURL, err := url.Parse(*urlFlag)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}
	startURL.Path = strings.TrimSuffix(*&startURL.Path, "/")

	return &Config{
		URL:      startURL,
		MaxDepth: *depth,
	}, nil
}
