package main

import (
	"database/sql"
	"yuki_buy_log/validators"
)

type Server struct {
	db        *sql.DB
	validator validators.Validator
	auth      *Authenticator
}

func NewServer(db *sql.DB, v validators.Validator, a *Authenticator) *Server {
	return &Server{db: db, validator: v, auth: a}
}
