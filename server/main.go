package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

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
	defer db.Close()

	auth := NewAuthenticator([]byte("secret"))
	srv := NewServer(db, NewValidator(), auth)

	mux := http.NewServeMux()
	mux.Handle("/products", auth.Middleware(http.HandlerFunc(srv.productsHandler)))
	mux.Handle("/purchases", auth.Middleware(http.HandlerFunc(srv.purchasesHandler)))
	mux.HandleFunc("/register", srv.registerHandler)
	mux.HandleFunc("/login", srv.loginHandler)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
