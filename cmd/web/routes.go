package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := &http.ServeMux{}

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", app.index)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/list", app.snippetList)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	return mux
}
