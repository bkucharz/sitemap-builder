package main

import (
	"flag"
	"net/url"
	"os"
	"strings"
	"testing"
)

func TestParseFlags(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	tests := []struct {
		name    string
		args    []string
		wantURL string
		wantErr bool
	}{
		{
			name:    "default values",
			args:    []string{"cmd"},
			wantURL: "https://go.dev",
		},
		{
			name:    "custom url",
			args:    []string{"cmd", "-url", "https://example.com/"},
			wantURL: "https://example.com",
		},
		{
			name:    "invalid url",
			args:    []string{"cmd", "-url", ":invalid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(tt.name, flag.ContinueOnError)
			os.Args = tt.args

			got, err := parseFlags()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFlags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			if got.URL.String() != tt.wantURL {
				t.Errorf("parseFlags() URL = %v, want %v", got.URL.String(), tt.wantURL)
			}
		})
	}
}

func TestXMLGeneration(t *testing.T) {
	conf := &Config{
		URL:      &url.URL{Scheme: "https", Host: "example.com"},
		MaxDepth: 1,
	}

	paths := []string{"/foo", "/bar"}
	xml := generateXML(paths, *conf)

	if !strings.Contains(xml, "<loc>https://example.com/foo</loc>") {
		t.Error("XML missing expected URL")
	}
	if !strings.Contains(xml, `xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"`) {
		t.Error("XML missing namespace")
	}
}
