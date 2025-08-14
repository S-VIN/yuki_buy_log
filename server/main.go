package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://user:pass@localhost:5432/yukibuylog?sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 30; i++ {
		if err = db.Ping(); err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	auth := NewAuthenticator([]byte("secret"))
	srv := NewServer(db, NewValidator(), auth)

	mux := http.NewServeMux()
	mux.Handle("/products", auth.Middleware(http.HandlerFunc(srv.productsHandler)))
	mux.Handle("/purchases", auth.Middleware(http.HandlerFunc(srv.purchasesHandler)))
	mux.Handle("/family/members", auth.Middleware(http.HandlerFunc(srv.familyMembersHandler)))
	mux.Handle("/family/invitations", auth.Middleware(http.HandlerFunc(srv.familyInvitationsHandler)))
	mux.Handle("/family/invite", auth.Middleware(http.HandlerFunc(srv.familyInviteHandler)))
	mux.Handle("/family/respond", auth.Middleware(http.HandlerFunc(srv.familyRespondHandler)))
	mux.Handle("/family/leave", auth.Middleware(http.HandlerFunc(srv.familyLeaveHandler)))
	mux.HandleFunc("/register", srv.registerHandler)
	mux.HandleFunc("/login", srv.loginHandler)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", enableCORS(mux)))
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
