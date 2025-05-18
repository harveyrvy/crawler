package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type Result struct {
	Url   string
	Links []string
}
type Crawler struct {
	startURL *url.URL
	visited  map[string]bool
	results  []Result
	client   *http.Client
	timeout  time.Duration
}

func NewCrawler(startURL string, timeoutInSeconds string) (*Crawler, error) {
	url, err := url.Parse(startURL)
	if err != nil {
		return nil, fmt.Errorf("Invalid url %v", err)
	}
	timeout, err := strconv.Atoi(timeoutInSeconds)
	if err != nil {
		return nil, fmt.Errorf("Invalid timeout %v", err)
	}
	crawler := &Crawler{
		startURL: url,
		client:   &http.Client{},
		visited:  make(map[string]bool),
		results:  []Result{},
		timeout:  time.Duration(timeout) * time.Second,
	}
	return crawler, nil
}

func (c *Crawler) crawl() error {
	startTime := time.Now()
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

	for len(uncheckedLinks) > 0 {
		if !c.visited[uncheckedLinks[0]] {
			newLinks, err := c.crawlPage(uncheckedLinks[0])
			if err != nil {
				return err
			}

			for _, l := range newLinks {
				if !c.visited[l] {
					uncheckedLinks = append(uncheckedLinks, l)
				}
			}
		}
		uncheckedLinks = uncheckedLinks[1:]
		if time.Now().After(startTime.Add(c.timeout)) {
			break
		}
	}
	return nil
}

func (c *Crawler) crawlPage(url string) ([]string, error) {
	fmt.Printf("crawling %v\n", url)
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("error getting doc %v", err)
		return []string{}, err
	}
	defer response.Body.Close()
	doc, err := html.Parse(response.Body)
	if err != nil {
		fmt.Printf("error parsing doc %v", err)
		return []string{}, err
	}
	c.visited[url] = true
	links, err := c.getAllLinks(doc)
	if err != nil {
		fmt.Printf("error getting links from page %v", err)
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
						fmt.Printf("Error parsing link %v", err)
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
