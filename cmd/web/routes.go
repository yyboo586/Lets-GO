package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.Dir("./ui/static"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static/", fileServer))

	router.HandlerFunc(http.MethodGet, "/", app.index)
	router.HandlerFunc(http.MethodGet, "/snippet/views/:id", app.snippetView)
	router.HandlerFunc(http.MethodGet, "/snippet/views", app.snippetList)
	router.HandlerFunc(http.MethodPost, "/snippet/views", app.snippetCreate)

	// Create the middleware chain as normal.
	standard := alice.New(app.panicRecover, app.logRequest)

	return standard.Then(router)
}
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		app.infoLogger.Printf("%s - %s %s %s %v\n", r.RemoteAddr, r.Proto, r.Method, r.RequestURI, time.Since(start))
	})
}

func (app *application) panicRecover(next http.Handler) http.Handler {
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
