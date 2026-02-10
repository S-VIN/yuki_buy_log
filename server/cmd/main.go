package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"yuki_buy_log/internal/auth"
	"yuki_buy_log/internal/handlers"
	"yuki_buy_log/internal/tasks"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	authenticator := auth.NewAuthenticator()

	mux := newServeMux(authenticator)
	srv := newHTTPServer(mux)
	scheduler := newScheduler()

	scheduler.Start()

	go func() {
		log.Printf("Server started on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")

	scheduler.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}

func newServeMux(authenticator *auth.Authenticator) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.Handle("/products", authenticator.Middleware(handlers.ProductsHandler(authenticator)))
	mux.Handle("/purchases", authenticator.Middleware(handlers.PurchasesHandler(authenticator)))
	mux.Handle("/group", authenticator.Middleware(handlers.GroupHandler(authenticator)))
	mux.Handle("/invite", authenticator.Middleware(handlers.InviteHandler(authenticator)))

	mux.HandleFunc("/register", handlers.RegisterHandler(authenticator))
	mux.HandleFunc("/login", handlers.LoginHandler(authenticator))

	return mux
}

func newHTTPServer(mux *http.ServeMux) *http.Server {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	corsOrigin := os.Getenv("CORS_ORIGIN")
	if corsOrigin == "" {
		corsOrigin = "http://localhost:5173"
	}

	return &http.Server{
		Addr:    ":" + port,
		Handler: enableCORS(mux, corsOrigin),
	}
}

func newScheduler() *tasks.Scheduler {
	scheduler := tasks.NewScheduler()
	scheduler.AddTask(tasks.Task{
		Name:     "cleanup_old_invites",
		Interval: 5 * time.Minute,
		Run:      tasks.CleanupOldInvites(),
	})
	return scheduler
}

func enableCORS(next http.Handler, corsOrigin string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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