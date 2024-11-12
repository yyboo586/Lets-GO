package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type config struct {
	Addr        string // HTTP Server Address
	LogFilePath string // Log file path
}

type application struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

func main() {
	config := &config{}
	flag.StringVar(&config.Addr, "addr", ":8080", "http server address")
	// flag.StringVar(&config.LogFilePath, "log", "/var/log/snippest/access.log", "log file path")
	flag.Parse()

	infoLogger := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLogger := log.New(os.Stderr, "\033[31mERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		infoLogger:  infoLogger,
		errorLogger: errorLogger,
	}

	srv := http.Server{
		Addr:     config.Addr,
		Handler:  app.routes(),
		ErrorLog: errorLogger,
	}

	infoLogger.Printf("Server listening at %s", config.Addr)
	if err := srv.ListenAndServe(); err != nil {
		errorLogger.Fatal(err)
	}
}
