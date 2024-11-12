package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

func index(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/base.html",
		"./ui/html/pages/home.html",
		"./ui/html/partials/nav.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("服务器内部错误: %v", err), http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("服务器内部错误: %v", err), http.StatusInternalServerError)
		return
	}
}

func create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain") // 避免内容嗅探
	w.Write([]byte("POST /create\n"))
}

func urlQueryParam(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	age, err := strconv.Atoi(r.URL.Query().Get("age"))
	if err != nil {
		http.Error(w, fmt.Sprintf("参数错误: %v", err), http.StatusBadRequest)
		return
	}

	log.Printf("name=%s, age=%d\n", name, age)
	w.WriteHeader(http.StatusOK)
}
