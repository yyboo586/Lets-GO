package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.ServeMux{}

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", index)
	mux.HandleFunc("/create", create)
	mux.HandleFunc("/url-query-param", urlQueryParam)

	log.Println("Server listening at :8080")
	if err := http.ListenAndServe(":8080", &mux); err != nil {
		panic(err)
	}
}
