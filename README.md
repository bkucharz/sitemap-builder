# 🗺️ Sitemap Builder

[![Go Version](https://img.shields.io/badge/go-1.23+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/bkucharz/sitemap-builder)](https://goreportcard.com/report/github.com/bkucharz/sitemap-builder)

A lightweight Go CLI tool that crawls websites and generates XML sitemaps.
## ✨ Features

- **Website Crawling** - Discovers all pages up to specified depth
- **Standard Compliant** - Generates valid XML sitemaps
- **Link Normalization** - Handles relative/absolute URLs
- **Duplicate Prevention** - Automatically filters duplicate pages

## 🚀 Quick Start

### Installation

```bash
go install github.com/bkucharz/sitemap-builder/cmd/sitemap-builder@latest
```

### Basic Usage

```bash
sitemap-builder -url https://example.com -depth 2
```

### Example Output

```xml
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
 <url>
  <loc>https://example.com/</loc>
 </url>
 <url>
  <loc>https://example.com/about</loc>
 </url>
</urlset>
```

## 📚 Command Line Options

| Flag     | Default          | Description                      |
|----------|------------------|----------------------------------|
| `-url`   | `https://go.dev` | Starting URL to crawl            |
| `-depth` | `1`              | Maximum depth of links to follow |

### Development Setup

```bash
git clone https://github.com/bkucharz/sitemap-builder
cd sitemap-builder
go build ./cmd/sitemap-builder
```

## 📜 License

[MIT](LICENSE)