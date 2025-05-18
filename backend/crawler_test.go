package main

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestGetAllLinks(t *testing.T) {
	c, err := NewCrawler("example.com", "5")
	if err != nil {
		t.Fatal("Error creating Crawler", err)
	}
	const exampleDoc = `
<html>
  <div>
    <a href="/link/1"></a>
    <a href="/link/2"></a>
  </div>
  <a href="/link/3"></a>
  <div>
    <div>
      <a href="/link/4"></a>
    </div>
    <a href="/link/5"></a>
  </div>
  <a href="/link/6"></a>
</html>`
	doc, err := html.Parse(strings.NewReader(exampleDoc))
	if err != nil {
		t.Fatal("Error parsing example testing document", err)
	}
	expected := []string{
		"https://example.com/link/1",
		"https://example.com/link/2",
		"https://example.com/link/3",
		"https://example.com/link/4",
		"https://example.com/link/5",
		"https://example.com/link/6",
	}

	got, err := c.getAllLinks(doc)
	if len(got) != len(expected) {
		t.Fatalf("Wrong amount of links. Expected %v, got %v, %v", len(expected), len(got), err)
	}
}
