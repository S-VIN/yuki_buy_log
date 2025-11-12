package database

import (
	"fmt"
	"log"
	"yuki_buy_log/models"

	"github.com/lib/pq"
)

func (d *DatabaseManager) GetAllPurchases() ([]models.Purchase, error) {
	rows, err := d.db.Query(`SELECT id, product_id, quantity, price, date, store, tags, receipt_id, user_id FROM purchases`)
	if err != nil {
		return nil, fmt.Errorf("failed to get all purchases: %w", err)
	}
	defer rows.Close()

	var purchases []models.Purchase
	for rows.Next() {
		var p models.Purchase
		err := rows.Scan(&p.Id, &p.ProductId, &p.Quantity, &p.Price, &p.Date, &p.Store, pq.Array(&p.Tags), &p.ReceiptId, &p.UserId)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		purchases = append(purchases, p)
	}
	return purchases, nil
}

func (d *DatabaseManager) GetPurchasesByUserIds(userIds []models.UserId) ([]models.Purchase, error) {
	if len(userIds) == 0 {
		return []models.Purchase{}, nil
	}

	rows, err := d.db.Query(`SELECT id, product_id, quantity, price, date, store, tags, receipt_id, user_id FROM purchases WHERE user_id = ANY($1)`, pq.Array(userIds))
	if err != nil {
		return nil, fmt.Errorf("failed to get purchases for users: %w", err)
	}
	defer rows.Close()

	var purchases []models.Purchase
	for rows.Next() {
		var p models.Purchase
		err := rows.Scan(&p.Id, &p.ProductId, &p.Quantity, &p.Price, &p.Date, &p.Store, pq.Array(&p.Tags), &p.ReceiptId, &p.UserId)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		purchases = append(purchases, p)
	}
	return purchases, nil
}

func (d *DatabaseManager) AddPurchase(purchase *models.Purchase) error {
	err := d.db.QueryRow(`INSERT INTO purchases (product_id, quantity, price, date, store, tags, receipt_id, user_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`,
		purchase.ProductId, purchase.Quantity, purchase.Price, purchase.Date, purchase.Store, pq.Array(purchase.Tags), purchase.ReceiptId, purchase.UserId).Scan(&purchase.Id)
	if err != nil {
		log.Printf("Failed to insert purchase: %v", err)
		return err
	}
	return nil
}

func (d *DatabaseManager) DeletePurchase(purchaseId models.PurchaseId, userId models.UserId) error {
	result, err := d.db.Exec(`DELETE FROM purchases WHERE id = $1 AND user_id = $2`, purchaseId, userId)
	if err != nil {
		log.Printf("Failed to delete purchase: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Failed to check rows affected: %v", err)
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("purchase with id %d not found for user %d", purchaseId, userId)
	}

	return nil
}
