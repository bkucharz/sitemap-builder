package link

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Link
		wantErr  bool
	}{
		{
			name:  "single link",
			input: `<a href="/foo">Foo</a>`,
			expected: []Link{
				{Href: "/foo", Text: "Foo"},
			},
		},
		{
			name:  "multiple links",
			input: `<a href="/foo">Foo</a><a href="/bar">Bar</a>`,
			expected: []Link{
				{Href: "/foo", Text: "Foo"},
				{Href: "/bar", Text: "Bar"},
			},
		},
		{
			name:     "invalid html",
			input:    `<a href="/foo"`,
			expected: []Link{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			got, err := Parse(r)

			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != len(tt.expected) {
				t.Errorf("got %d links, want %d", len(got), len(tt.expected))
			}

			for i, link := range got {
				if link.Href != tt.expected[i].Href || link.Text != tt.expected[i].Text {
					t.Errorf("link %d: got %+v, want %+v", i, link, tt.expected[i])
				}
			}
		})
	}
}
