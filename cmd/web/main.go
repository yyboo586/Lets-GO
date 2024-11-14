package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"snippetbox/internal/models"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
)

type config struct {
	Addr        string // HTTP Server Address
	LogFilePath string // Log file path

	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
}

type application struct {
	infoLogger     *log.Logger
	errorLogger    *log.Logger
	snippets       *models.SnippetModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	config := &config{}
	flag.StringVar(&config.Addr, "addr", ":8080", "http server address")
	// flag.StringVar(&config.LogFilePath, "log", "/var/log/snippest/access.log", "log file path")
	flag.StringVar(&config.DBHost, "dbhost", "localhost", "database host")
	flag.IntVar(&config.DBPort, "dbport", 3306, "database port")
	flag.StringVar(&config.DBUser, "dbuser", "root", "database user")
	flag.StringVar(&config.DBPassword, "dbpass", "12345678", "database password")
	flag.StringVar(&config.DBName, "dbname", "snippetbox", "database name")
	flag.Parse()

	infoLogger := log.New(os.Stdout, "\033[0mINFO\t", log.Ldate|log.Ltime)
	errorLogger := log.New(os.Stderr, "\033[31mERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
	dbPool, err := openDB(dsn)
	if err != nil {
		errorLogger.Fatal(err)
	}
	defer dbPool.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLogger.Fatal(err)
	}

	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Store = mysqlstore.New(dbPool)
	sessionManager.Cookie.Secure = true

	app := &application{
		infoLogger:     infoLogger,
		errorLogger:    errorLogger,
		snippets:       &models.SnippetModel{DB: dbPool},
		users:          &models.UserModel{DB: dbPool},
		templateCache:  templateCache,
		formDecoder:    form.NewDecoder(),
		sessionManager: sessionManager,
	}

	// tlsConfig := &tls.Config{
	// 	CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	// }

	srv := http.Server{
		Addr:     config.Addr,
		Handler:  app.routes(),
		ErrorLog: errorLogger,
		// TLSConfig: tlsConfig,

		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLogger.Printf("Server listening at %s", config.Addr)
	// if err := srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem"); err != nil {
	// 	errorLogger.Fatal(err)
	// }
	if err := srv.ListenAndServe(); err != nil {
		errorLogger.Fatal(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
