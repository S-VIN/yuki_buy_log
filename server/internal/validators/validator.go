package validators

import (
	"errors"
	"regexp"

	"yuki_buy_log/internal/domain"
)

var (
	// reValidName allows Unicode letters, digits, and spaces
	reValidName = regexp.MustCompile(`^[\p{L}\p{N}\s]+$`)
)

// ValidateProduct validates a product.
func ValidateProduct(p *domain.Product) error {
	if len(p.Name) == 0 || len(p.Name) > 30 || !reValidName.MatchString(p.Name) {
		return errors.New("invalid name")
	}
	if len(p.Volume) == 0 || len(p.Volume) > 10 {
		return errors.New("invalid volume")
	}
	if len(p.Brand) == 0 || len(p.Brand) > 30 || !reValidName.MatchString(p.Brand) {
		return errors.New("invalid brand")
	}
	if len(p.DefaultTags) > 10 {
		return errors.New("too many default tags")
	}
	for _, tag := range p.DefaultTags {
		if len(tag) == 0 || len(tag) > 20 || !reValidName.MatchString(tag) {
			return errors.New("invalid default tag")
		}
	}
	return nil
}

// ValidatePurchase validates a purchase.
func ValidatePurchase(p *domain.Purchase) error {
	if p.ProductId <= 0 {
		return errors.New("invalid product_id")
	}
	if p.Quantity < 1 || p.Quantity > 100000 {
		return errors.New("invalid quantity")
	}
	if p.Price < 1 || p.Price > 100000000 {
		return errors.New("invalid price")
	}
	if p.Date.IsZero() {
		return errors.New("invalid date")
	}
	if len(p.Store) == 0 || len(p.Store) > 30 || !reValidName.MatchString(p.Store) {
		return errors.New("invalid store")
	}
	if len(p.Tags) > 10 {
		return errors.New("too many tags")
	}
	for _, tag := range p.Tags {
		if len(tag) == 0 || len(tag) > 20 || !reValidName.MatchString(tag) {
			return errors.New("invalid tag")
		}
	}
	return nil
}
