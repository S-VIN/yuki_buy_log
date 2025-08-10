//go:build integration

package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type productsResp struct {
	Products []Product `json:"products"`
}

func setupPostgres(t *testing.T) (*sql.DB, func()) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		Env:          map[string]string{"POSTGRES_PASSWORD": "pass", "POSTGRES_USER": "user", "POSTGRES_DB": "yukibuylog"},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
	}
	pg, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{ContainerRequest: req, Started: true})
	if err != nil {
		t.Skipf("docker not available: %v", err)
	}
	host, err := pg.Host(ctx)
	if err != nil {
		t.Fatalf("host: %v", err)
	}
	port, err := pg.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("port: %v", err)
	}
	dsn := fmt.Sprintf("postgres://user:pass@%s:%s/yukibuylog?sslmode=disable", host, port.Port())
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	// wait for db ready
	for i := 0; i < 20; i++ {
		if err := db.Ping(); err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	schema, err := os.ReadFile(filepath.Join("..", "postgres", "schema.sql"))
	if err != nil {
		t.Fatalf("schema: %v", err)
	}
	if _, err := db.Exec(string(schema)); err != nil {
		t.Fatalf("exec schema: %v", err)
	}
	data, err := os.ReadFile(filepath.Join("..", "postgres", "testdata.sql"))
	if err != nil {
		t.Fatalf("data: %v", err)
	}
	if _, err := db.Exec(string(data)); err != nil {
		t.Fatalf("exec data: %v", err)
	}
	return db, func() {
		db.Close()
		pg.Terminate(ctx)
	}
}

func startServer(db *sql.DB) *httptest.Server {
	srv := NewServer(db, NewValidator())
	mux := http.NewServeMux()
	mux.HandleFunc("/products", srv.productsHandler)
	mux.HandleFunc("/purchases", srv.purchasesHandler)
	return httptest.NewServer(mux)
}

func TestSystemEndpoints(t *testing.T) {
	db, teardown := setupPostgres(t)
	defer teardown()
	ts := startServer(db)
	defer ts.Close()

	// initial products
	resp, err := http.Get(ts.URL + "/products")
	if err != nil {
		t.Fatalf("get products: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var pr productsResp
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(pr.Products) == 0 {
		t.Fatalf("expected products")
	}

	// create product
	body := strings.NewReader(`{"name":"Water","volume":"1l","brand":"Brand3","category":"Drink","description":""}`)
	resp, err = http.Post(ts.URL+"/products", "application/json", body)
	if err != nil {
		t.Fatalf("post product: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	// create purchase
	body = strings.NewReader(`{"product_id":1,"quantity":1,"price":50,"date":"2023-03-10","store":"Store","receipt_id":1}`)
	resp, err = http.Post(ts.URL+"/purchases", "application/json", body)
	if err != nil {
		t.Fatalf("post purchase: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	resp.Body.Close()
}
