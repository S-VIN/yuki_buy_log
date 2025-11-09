package main

import (
	"database/sql"
)

type Server struct {
	db   *sql.DB
	auth *Authenticator
}

func NewServer(db *sql.DB, a *Authenticator) *Server {
	return &Server{db: db, auth: a}
}
