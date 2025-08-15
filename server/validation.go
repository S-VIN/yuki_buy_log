package main

import (
	"errors"
	"regexp"
	"time"
)

// Validator validates incoming products and purchases.
type Validator interface {
	ValidateProduct(*Product) error
	ValidatePurchase(*Purchase) error
	ValidateLogin(string) error
}

type validator struct{}

var (
	reLetters       = regexp.MustCompile(`^[A-Za-z]+$`)
	reLettersDigits = regexp.MustCompile(`^[A-Za-z0-9]+$`)
)

// NewValidator returns a Validator implementation.
func NewValidator() Validator { return validator{} }

func (validator) ValidateProduct(p *Product) error {
	if len(p.Name) == 0 || len(p.Name) > 30 || !reLetters.MatchString(p.Name) {
		return errors.New("invalid name")
	}
	if len(p.Volume) == 0 || len(p.Volume) > 10 {
		return errors.New("invalid volume")
	}
	if len(p.Brand) == 0 || len(p.Brand) > 30 || !reLettersDigits.MatchString(p.Brand) {
		return errors.New("invalid brand")
	}
	if len(p.Category) == 0 || len(p.Category) > 30 || !reLetters.MatchString(p.Category) {
		return errors.New("invalid category")
	}
	if len(p.Description) > 150 {
		return errors.New("invalid description")
	}
	return nil
}

func (validator) ValidatePurchase(p *Purchase) error {
	if p.ProductId <= 0 {
		return errors.New("invalid product_id")
	}
	if p.Quantity < 1 || p.Quantity > 100000 {
		return errors.New("invalid quantity")
	}
	if p.Price < 1 || p.Price > 100000000 {
		return errors.New("invalid price")
	}
	if p.Date == "" {
		return errors.New("invalid date")
	}
	if _, err := time.Parse("2006-01-02", p.Date); err != nil {
		return errors.New("invalid date")
	}
	if len(p.Store) == 0 || len(p.Store) > 30 || !reLetters.MatchString(p.Store) {
		return errors.New("invalid store")
	}
	return nil
}

func (validator) ValidateLogin(l string) error {
	if len(l) == 0 || len(l) > 50 || !reLettersDigits.MatchString(l) {
		return errors.New("invalid login")
	}
	return nil
}
