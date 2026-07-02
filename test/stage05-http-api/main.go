package main

import (
	"fmt"
	"net/http"
)

type Server struct {
	Name   string `json:"name"`
	IP     string `json:"ip"`
	Status string `json:"status"`
}

func main() {

	var servers = []Server{
		{"web-01", "10.0.0.5", "running"},
		{"db-01", "10.0.0.21", "running"},
		{"cache-01", "10.0.0.30", "stopped"},
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/server", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

	})
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	fmt.Println("http://localhost:8000/health")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
