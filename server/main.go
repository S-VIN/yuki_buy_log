package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"yuki_buy_log/database"
	"yuki_buy_log/handlers"
	"yuki_buy_log/tasks"

	_ "github.com/lib/pq"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("Initializing authenticator...")
	auth := NewAuthenticator([]byte("secret")) // TODO change key

	log.Println("Initializing database...")
	handlers.DB = database.NewPostgresDB()

	log.Println("Setting up HTTP routes...")

	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/yuki_buy_log?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}
	// We need db for potential future use, so we keep it
	_ = db

	mux := http.NewServeMux()

	// Health check endpoint (no auth required)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.Handle("/products", auth.Middleware(handlers.ProductsHandler(auth)))
	mux.Handle("/purchases", auth.Middleware(handlers.PurchasesHandler(auth)))
	mux.Handle("/group", auth.Middleware(handlers.GroupHandler(auth)))
	mux.Handle("/invite", auth.Middleware(handlers.InviteHandler(auth)))
	mux.HandleFunc("/register", handlers.RegisterHandler(auth))
	mux.HandleFunc("/login", handlers.LoginHandler(auth))

	// Setup and start scheduler
	log.Println("Setting up scheduler...")
	scheduler := tasks.NewScheduler()
	scheduler.AddTask(tasks.Task{
		Name:     "cleanup_old_invites",
		Interval: 5 * time.Minute,
		Run:      tasks.CleanupOldInvites(),
	})
	scheduler.Start()
	defer scheduler.Stop()

	// Get server port from environment
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	// Setup graceful shutdown
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: enableCORS(mux),
	}

	// Channel to listen for interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		log.Printf("Server started on :%s", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-stop
	log.Println("Shutting down gracefully...")
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		corsOrigin := os.Getenv("CORS_ORIGIN")
		if corsOrigin == "" {
			corsOrigin = "http://localhost:5173"
		}
		log.Printf("CORS middleware processing request: %s %s from %s", r.Method, r.URL.Path, r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Origin", corsOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			log.Printf("Handling OPTIONS preflight request for %s", r.URL.Path)
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
