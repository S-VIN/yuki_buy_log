package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"yuki_buy_log/handlers"
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

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", enableCORS(mux)))
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
