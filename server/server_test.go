package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

type mockValidator struct {
	productErr  error
	purchaseErr error
}

func (m *mockValidator) ValidateProduct(*Product) error   { return m.productErr }
func (m *mockValidator) ValidatePurchase(*Purchase) error { return m.purchaseErr }

func newTestServer(t *testing.T) (*Server, sqlmock.Sqlmock, *mockValidator, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock: %v", err)
	}
	v := &mockValidator{}
	srv := NewServer(db, v, NewAuthenticator([]byte("secret")))
	return srv, mock, v, func() { db.Close() }
}

func TestGetProducts(t *testing.T) {
	srv, mock, _, close := newTestServer(t)
	defer close()

	fixed := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	rows := sqlmock.NewRows([]string{"id", "name", "volume", "brand", "category", "description", "creation_date"}).
		AddRow(1, "Tea", "500ml", "Brand1", "Drink", "Green tea", fixed)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, volume, brand, category, description, creation_date FROM products")).
		WillReturnRows(rows)

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()
	srv.productsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestCreateProduct(t *testing.T) {
	srv, mock, _, close := newTestServer(t)
	defer close()

	body := `{"name":"Tea","volume":"500ml","brand":"Brand1","category":"Drink","description":"Green tea"}`

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO products (name, volume, brand, category, description, creation_date) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id")).
		WithArgs("Tea", "500ml", "Brand1", "Drink", "Green tea", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(body))
	w := httptest.NewRecorder()
	srv.productsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if strings.Contains(w.Body.String(), "creation_date") {
		t.Fatalf("creation_date should not be in response")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetPurchases(t *testing.T) {
	srv, mock, _, close := newTestServer(t)
	defer close()

	fixed := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)
	rows := sqlmock.NewRows([]string{"id", "product_id", "quantity", "price", "date", "store", "receipt_id"}).
		AddRow(1, 1, 2, 100, fixed, "Store", 1)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, product_id, quantity, price, date, store, receipt_id FROM purchases WHERE user_id=$1")).
		WithArgs(int64(1)).
		WillReturnRows(rows)

	req := httptest.NewRequest(http.MethodGet, "/purchases", nil)
	req = req.WithContext(context.WithValue(req.Context(), UserIDKey, int64(1)))
	w := httptest.NewRecorder()
	srv.purchasesHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestCreatePurchase(t *testing.T) {
	srv, mock, _, close := newTestServer(t)
	defer close()

	body := `{"product_id":1,"quantity":2,"price":100,"date":"2023-03-01","store":"Store","receipt_id":1}`

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO purchases (product_id, quantity, price, date, store, receipt_id, user_id) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id")).
		WithArgs(int64(1), 2, 100, sqlmock.AnyArg(), "Store", int64(1), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	req := httptest.NewRequest(http.MethodPost, "/purchases", strings.NewReader(body))
	req = req.WithContext(context.WithValue(req.Context(), UserIDKey, int64(1)))
	w := httptest.NewRecorder()
	srv.purchasesHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestCreateProductInvalid(t *testing.T) {
	srv, _, v, close := newTestServer(t)
	defer close()

	body := `{"name":"Tea","volume":"500ml","brand":"Brand1","category":"Drink","description":"Green"}`
	v.productErr = errors.New("invalid")

	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(body))
	w := httptest.NewRecorder()
	srv.productsHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCreatePurchaseInvalidStore(t *testing.T) {
	srv, _, v, close := newTestServer(t)
	defer close()

	body := `{"product_id":1,"quantity":2,"price":100,"date":"2023-03-01","store":"Store","receipt_id":1}`
	v.purchaseErr = errors.New("invalid store")

	req := httptest.NewRequest(http.MethodPost, "/purchases", strings.NewReader(body))
	w := httptest.NewRecorder()
	srv.purchasesHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCreatePurchaseMissingDate(t *testing.T) {
	srv, _, v, close := newTestServer(t)
	defer close()

	body := `{"product_id":1,"quantity":2,"price":100,"store":"Store","receipt_id":1}`
	v.purchaseErr = errors.New("invalid date")

	req := httptest.NewRequest(http.MethodPost, "/purchases", strings.NewReader(body))
	w := httptest.NewRecorder()
	srv.purchasesHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
