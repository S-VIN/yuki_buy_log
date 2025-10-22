package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"yuki_buy_log/handlers"
	"yuki_buy_log/scheduler"
	"yuki_buy_log/tasks"
	"yuki_buy_log/validators"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/yuki_buy_log?sslmode=disable"
	}

	log.Printf("Connecting to database with DSN: %s", dsn)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}

	log.Println("Waiting for database connection...")
	for i := 0; i < 30; i++ {
		if err = db.Ping(); err == nil {
			log.Println("Database connection established successfully")
			break
		}
		log.Printf("Database connection attempt %d failed, retrying in 1 second...", i+1)
		time.Sleep(time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to connect to database after 30 attempts: %v", err)
	}
	defer db.Close()

	log.Println("Initializing authenticator...")
	auth := NewAuthenticator([]byte("secret"))

	log.Println("Setting up HTTP routes...")

	// Создаем зависимости для handlers
	deps := &handlers.Dependencies{
		DB:        db,
		Validator: validators.NewValidator(),
		Auth:      auth,
	}

	mux := http.NewServeMux()
	mux.Handle("/products", auth.Middleware(handlers.ProductsHandler(deps)))
	mux.Handle("/purchases", auth.Middleware(handlers.PurchasesHandler(deps)))
	mux.Handle("/group", auth.Middleware(handlers.GroupHandler(deps)))
	mux.Handle("/invite", auth.Middleware(handlers.InviteHandler(deps)))
	mux.HandleFunc("/register", handlers.RegisterHandler(deps))
	mux.HandleFunc("/login", handlers.LoginHandler(deps))

	// Setup and start scheduler
	log.Println("Setting up scheduler...")
	sched := scheduler.New()
	sched.AddTask(scheduler.Task{
		Name:     "cleanup_old_invites",
		Interval: 5 * time.Minute,
		Run:      tasks.CleanupOldInvites(db),
	})
	sched.Start()
	defer sched.Stop()

	// Setup graceful shutdown
	srv := &http.Server{
		Addr:    ":8080",
		Handler: enableCORS(mux),
	}

	// Channel to listen for interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		log.Println("Server started on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-stop
	log.Println("Shutting down gracefully...")
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("CORS middleware processing request: %s %s from %s", r.Method, r.URL.Path, r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
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
