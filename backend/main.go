package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Success bool     `json:"success"`
	Results []Result `json:"results"`
}

func getCrawl(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Need options for CORS
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")

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

func handle() {
	http.HandleFunc("/", getCrawl)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	handle()
}
