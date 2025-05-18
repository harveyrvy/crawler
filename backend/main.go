package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Response struct {
	Success bool     `json:success`
	Results []Result `json:results`
}

func getCrawl(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	timeout := r.URL.Query().Get("timeout")
	if timeout == "" {
		timeout = "5"
	}
	c, err := NewCrawler(url, timeout)
	if err != nil {
		response := Response{
			Success: false,
			Results: nil,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	err = c.crawl()
	if err != nil {
		response := Response{
			Success: false,
			Results: nil,
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	response := Response{
		Success: true,
		Results: c.results,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	return
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "homepage")
}

func handle() {
	http.HandleFunc("/", getCrawl)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	handle()
}
