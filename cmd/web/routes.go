package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/justinas/nosurf"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.Dir("./ui/static"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static/", fileServer))

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.snippetList))

	router.Handler(http.MethodGet, "/snippet/views/:id", dynamic.ThenFunc(app.snippetView))

	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLogin))

	procted := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/snippet/views", procted.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/views", procted.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/user/logout", procted.ThenFunc(app.userLogoutPost))

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

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.IsAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}

func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	})

	return csrfHandler
}
