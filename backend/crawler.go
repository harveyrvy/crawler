package main

import (
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

type Result struct {
	url   string
	links []string
}

type Crawler struct {
	startURl *url.URL
	visited  map[string]bool
	results  []Result
	client   *http.Client
}

func NewCrawler(startURL string) (*Crawler, error) {
	url, err := url.Parse(startURL)
	if err != nil {
		return nil, fmt.Errorf("Invalid url %v", err)
	}
	crawler := &Crawler{
		startURl: url,
		client:   &http.Client{},
		visited:  make(map[string]bool),
		results:  []Result{},
	}
	return crawler, nil
}

func (c *Crawler) crawl() {
	c.crawlPage(c.startURl.String())
}

func (c *Crawler) crawlPage(url string) {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("error %v", err)
		return
	}
	defer response.Body.Close()
	// body, err := io.ReadAll(response.Body)
	// if err != nil {
	// 	fmt.Println("error %v", err)
	// 	return
	// }
	doc, err := html.Parse(response.Body)
	fmt.Println(doc)
	c.visited[url] = true
	c.getAllLinks(doc)
}

func (c *Crawler) getAllLinks(doc *html.Node) {
	fmt.Println(doc.FirstChild)

}

func main() {
	startURL := "https://bbc.com"
	c, err := NewCrawler(startURL)
	if err != nil {
		fmt.Println("error %v", err)
		return
	}
	c.crawl()
}
