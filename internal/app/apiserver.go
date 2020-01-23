package apiserver

import (
	"database/sql"
	"github.com/Andronovdima/tpark-db-forum/store"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
)

func Start() error {
	config := NewConfig()

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = ":5000"
	} else {
		port = ":" + port
	}
	config.BindAddr = port

	url :=  os.Getenv("DATABASE_URL")
	if len(url) != 0 {
		config.DatabaseURL = url
	}

	zapLogger, _ := zap.NewProduction()
	defer func() {
		if err := zapLogger.Sync(); err != nil {
			log.Println("HEHEHEG",err)
		}
	}()
	sugaredLogger := zapLogger.Sugar()

	srv, err := NewServer(config, sugaredLogger)
	if err != nil {
		return err
	}

	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Println(err)
		}
	}()

	srv.ConfigureServer(db)
	return http.ListenAndServe(config.BindAddr, srv)
}

func newDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(100)
	if err := store.CreateTables(db); err != nil {
		return nil, err
	}
	return db, nil
}
