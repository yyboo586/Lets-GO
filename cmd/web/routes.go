package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func (app *application) routes() http.Handler {
	mux := &http.ServeMux{}

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", app.index)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/list", app.snippetList)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	return panicRecover(secureHeaders(mux))
}
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		log.Printf("%s - %s %s %s %v\n", r.RemoteAddr, r.Proto, r.Method, r.RequestURI, time.Since(start))
	})
}

func panicRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
