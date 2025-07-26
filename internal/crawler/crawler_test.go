package crawler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"sitemap/internal/link"
	"testing"
)

func TestNormalizeLinks(t *testing.T) {
	base, _ := url.Parse("https://example.com")
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "relative link",
			input:    "/foo",
			expected: "https://example.com/foo",
		},
		{
			name:     "absolute link",
			input:    "https://example.com/bar",
			expected: "https://example.com/bar",
		},
		{
			name:     "link with fragment",
			input:    "/baz#section",
			expected: "https://example.com/baz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			links := []link.Link{{Href: tt.input}}
			urls := normalizeLinks(links, base)

			if len(urls) != 1 {
				t.Fatalf("expected 1 URL, got %d", len(urls))
			}

			if urls[0].String() != tt.expected {
				t.Errorf("got %s, want %s", urls[0].String(), tt.expected)
			}
		})
	}
}

func TestWithSameHost(t *testing.T) {
	base, _ := url.Parse("https://example.com")
	fn := withSameHost(base)

	tests := []struct {
		url  string
		want bool
	}{
		{"https://example.com/foo", true},
		{"https://sub.example.com/bar", false},
		{"http://example.com/baz", true},
		{"https://google.com/search", false},
	}

	for _, tt := range tests {
		u, _ := url.Parse(tt.url)
		got := fn(u)
		if got != tt.want {
			t.Errorf("%s: got %v, want %v", tt.url, got, tt.want)
		}
	}
}

func TestFetch(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html><body>
            <a href="/foo">Foo</a>
            <a href="/bar">Bar</a>  <!-- Changed from absolute to relative -->
        </body></html>`))
	}))
	defer ts.Close()

	testURL, _ := url.Parse(ts.URL)
	urls, err := Fetch(testURL)
	if err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}

	expected := []string{
		ts.URL + "/foo",
		ts.URL + "/bar",
	}

	if len(urls) != len(expected) {
		t.Fatalf("got %d URLs, want %d", len(urls), len(expected))
	}

	for i, u := range urls {
		if u.String() != expected[i] {
			t.Errorf("URL %d: got %s, want %s", i, u.String(), expected[i])
		}
	}
}
