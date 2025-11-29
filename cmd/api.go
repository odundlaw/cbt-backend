package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
)

type Application struct {
	config Config
	conn   *pgx.Conn
}

type Config struct {
	add string
	db  DBConfig
}

type DBConfig struct {
	dsn string // username and password for the
}

func (app *Application) mount() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all is good"))
	})

	// other routes
	//
	//
	//
	return r
}

func (app *Application) run(h http.Handler) error {
	server := &http.Server{
		Addr:         app.config.add,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server has started at add: %v", server.Addr)

	return server.ListenAndServe()
}
