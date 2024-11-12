package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := &http.ServeMux{}

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", app.index)
	mux.HandleFunc("/create", app.create)
	mux.HandleFunc("/url-query-param", app.urlQueryParam)

	return mux
}
