package database

import (
	"fmt"
	"log"
	"yuki_buy_log/internal/domain"

	"github.com/lib/pq"
)

func (d *DatabaseManager) GetAllProducts() ([]domain.Product, error) {
	rows, err := d.db.Query(`SELECT id, name, volume, brand, default_tags, user_id FROM products`)
	if err != nil {
		return nil, fmt.Errorf("failed to get all products: %w", err)
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		err := rows.Scan(&p.Id, &p.Name, &p.Volume, &p.Brand, pq.Array(&p.DefaultTags), &p.UserId)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		products = append(products, p)
	}
	return products, nil
}

func (d *DatabaseManager) GetProductById(id domain.ProductId) (*domain.Product, error) {
	var p domain.Product
	err := d.db.QueryRow(`SELECT id, name, volume, brand, default_tags, user_id FROM products WHERE id = $1`, id).
		Scan(&p.Id, &p.Name, &p.Volume, &p.Brand, pq.Array(&p.DefaultTags), &p.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to find product with id %d: %w", id, err)
	}
	return &p, nil
}

func (d *DatabaseManager) CreateProduct(product *domain.Product) error {
	err := d.db.QueryRow(`INSERT INTO products (name, volume, brand, default_tags, user_id) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		product.Name, product.Volume, product.Brand, pq.Array(product.DefaultTags), product.UserId).Scan(&product.Id)
	if err != nil {
		log.Printf("Failed to insert product: %v", err)
		return err
	}
	return nil
}

func (d *DatabaseManager) UpdateProduct(product *domain.Product) error {
	result, err := d.db.Exec(`UPDATE products SET name=$1, volume=$2, brand=$3, default_tags=$4 WHERE id=$5 AND user_id=$6`,
		product.Name, product.Volume, product.Brand, pq.Array(product.DefaultTags), product.Id, product.UserId)
	if err != nil {
		log.Printf("Failed to update product: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Failed to check rows affected: %v", err)
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product with id %d not found for user %d", product.Id, product.UserId)
	}

	return nil
}

func (d *DatabaseManager) DeleteProduct(id domain.ProductId, userId domain.UserId) error {
	result, err := d.db.Exec(`DELETE FROM products WHERE id = $1 AND user_id = $2`, id, userId)
	if err != nil {
		log.Printf("Failed to delete product: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Failed to check rows affected: %v", err)
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product with id %d not found for user %d", id, userId)
	}

	return nil
}
