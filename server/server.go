package main

import "database/sql"

type Server struct {
	db        *sql.DB
	validator Validator
	auth      *Authenticator
}

func NewServer(db *sql.DB, v Validator, a *Authenticator) *Server {
	return &Server{db: db, validator: v, auth: a}
}
