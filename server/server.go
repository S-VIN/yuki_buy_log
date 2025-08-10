package main

import "database/sql"

type Server struct {
	db        *sql.DB
	validator Validator
}

func NewServer(db *sql.DB, v Validator) *Server {
	return &Server{db: db, validator: v}
}
