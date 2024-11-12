package main

import (
	"net/http"
	"strconv"
	"text/template"
)

func (app *application) index(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/base.html",
		"./ui/html/pages/home.html",
		"./ui/html/partials/nav.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *application) create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain") // 避免内容嗅探
	w.Write([]byte("POST /create\n"))
}

func (app *application) urlQueryParam(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	age, err := strconv.Atoi(r.URL.Query().Get("age"))
	if err != nil || age < 0 {
		app.notFound(w)
		return
	}

	app.infoLogger.Printf("name=%s, age=%d\n", name, age)
	w.WriteHeader(http.StatusOK)
}
