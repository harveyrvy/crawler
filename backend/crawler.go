package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type Result struct {
	url   string
	links []string
}

type Crawler struct {
	startURL *url.URL
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
		startURL: url,
		client:   &http.Client{},
		visited:  make(map[string]bool),
		results:  []Result{},
	}
	return crawler, nil
}

func (c *Crawler) crawl() error {
	links, err := c.crawlPage(c.startURL.String())
	if err != nil {
		return err
	}
	fmt.Println("root page crawled")
	uncheckedLinks := []string{}
	for _, l := range links {
		if !c.visited[l] {
			uncheckedLinks = append(uncheckedLinks, l)
		}
	}

	for _, l := range uncheckedLinks {
		newLinks, err := c.crawlPage(l)
		if err != nil {
			return err
		}
		for _, l := range newLinks {
			if !c.visited[l] {
				uncheckedLinks = append(uncheckedLinks, l)
			}
		}

	}
	for _, r := range c.results {
		fmt.Println(r.url)
	}
	return nil
}

func (c *Crawler) crawlPage(url string) ([]string, error) {
	fmt.Printf("crawling %v\n", url)
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("error getting doc %v", err)
		return []string{}, err
	}
	defer response.Body.Close()
	doc, err := html.Parse(response.Body)
	if err != nil {
		fmt.Println("error parsing doc %v", err)
		return []string{}, err
	}
	c.visited[url] = true
	links, err := c.getAllLinks(doc)
	if err != nil {
		fmt.Println("error getting links from page %v", err)
		return []string{}, err
	}
	c.results = append(c.results, Result{url, links})
	return links, nil
}

func (c *Crawler) getAllLinks(doc *html.Node) ([]string, error) {
	var links []string
	uniqueLinks := make(map[string]bool)
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					u, err := url.Parse(strings.TrimSpace(attr.Val))
					if err != nil {
						fmt.Println("Error parsing link %v", err)
					}
					if (u.Host == c.startURL.Host) || u.Host == "" {
						resolvedUrl := c.startURL.ResolveReference(u)
						if !uniqueLinks[resolvedUrl.String()] {
							uniqueLinks[resolvedUrl.String()] = true
							links = append(links, resolvedUrl.String())
						}
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(doc)
	return links, nil

}

func main() {
	startURL := "https://www.bbc.co.uk"
	c, err := NewCrawler(startURL)
	if err != nil {
		fmt.Println("error %v", err)
		return
	}
	c.crawl()
}
