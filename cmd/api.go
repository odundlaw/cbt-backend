package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5"
	repo "github.com/odundlaw/cbt-backend/internal/adapters/postgresql/sqlc"
	"github.com/odundlaw/cbt-backend/internal/middlewares"
	"github.com/odundlaw/cbt-backend/internal/store"
	"github.com/odundlaw/cbt-backend/internal/users"
	"github.com/redis/go-redis/v9"
)

type Application struct {
	config Config
	conn   *pgx.Conn
	rdb    *redis.Client
}

type Config struct {
	add   string
	db    DBConfig
	redis RedisConfig
}

type DBConfig struct {
	dsn string // username and password for the
}

type RedisConfig struct {
	addr string // reids host localhost
}

func (app *Application) mount() http.Handler {
	r := chi.NewRouter()

	rdb := store.NewRedis(app.config.redis.addr)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all is good"))
	})

	userSerice := users.NewService(repo.New(app.conn))
	userHandler := users.NewHandler(userSerice, rdb)

	r.Mount("/", AuthRoutes(userHandler, rdb))

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

func AuthRoutes(handler *users.Handler, rdb *store.Redis) http.Handler {
	r := chi.NewRouter()

	// ——— USER AUTH ———
	r.Route("/api/auth", func(auth chi.Router) {
		// Public routes
		auth.Post("/register", handler.RegisterUser)
		auth.Post("/login", handler.LoginUser)
		auth.Post("/forgot-password", handler.ForgotPassword)

		// Protected routes
		auth.Group(func(protected chi.Router) {
			protected.Use(middlewares.AuthMiddleware(rdb))
			protected.Get("/refresh", handler.RefreshToken)
			protected.Post("/logout", handler.Logout)
		})
	})

	// ——— ADMIN AUTH ———
	r.Route("/api/admin", func(admin chi.Router) {
		// Public admin routes
		admin.Post("/register", handler.RegisterAdmin)
		admin.Post("/login", handler.LoginAdmin)
		admin.Post("/forgot-password", handler.AdminForgotPassword)

		// Protected admin routes
		admin.Group(func(protected chi.Router) {
			protected.Use(middlewares.AuthMiddleware(rdb))
			protected.Get("/refresh", handler.RefreshToken)
			protected.Post("/logout", handler.Logout)
		})
	})

	return r
}
