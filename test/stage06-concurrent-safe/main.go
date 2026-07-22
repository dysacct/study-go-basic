package main

import (
	"fmt"
	"net/http"

	"project-001/handler"
	"project-001/store"
)

func main() {

	serverStore := store.New()

	h := handler.New(serverStore)

	mux := http.NewServeMux()

	mux.HandleFunc("/servers", h.Servers)
	mux.HandleFunc("/server", h.Server)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	})

	fmt.Println("Server Start :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}