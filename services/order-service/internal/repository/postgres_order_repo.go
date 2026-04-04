package repository

import (
	"context"
	"database/sql"

	"github.com/Bharat1Rajput/flowpay/services/order-service/internal/model"
	"github.com/google/uuid"
)

type PostgresOrderRepo struct {
	db *sql.DB
}

func NewPostgresOrderRepo(db *sql.DB) *PostgresOrderRepo {
	return &PostgresOrderRepo{db: db}
}

func (r *PostgresOrderRepo) CreateOrder(ctx context.Context, order *model.Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert order
	_, err = tx.ExecContext(ctx, `
		INSERT INTO orders (id, user_id, status, total_amount, currency, delivery_addr, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`,
		order.ID,
		order.UserID,
		order.Status,
		order.TotalAmount,
		order.Currency,
		order.DeliveryAddr,
		order.Notes,
	)
	if err != nil {
		return err
	}

	// Insert items
	for _, item := range order.Items {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO order_items (id, order_id, item_name, quantity, unit_price, total_price)
			VALUES ($1, $2, $3, $4, $5, $6)
		`,
			item.ID,
			order.ID,
			item.ItemName,
			item.Quantity,
			item.UnitPrice,
			item.TotalPrice,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PostgresOrderRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Order, error) {

	row := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, status, total_amount, currency, delivery_addr, notes, created_at, updated_at
		FROM orders WHERE id = $1
	`, id)

	var order model.Order

	err := row.Scan(
		&order.ID,
		&order.UserID,
		&order.Status,
		&order.TotalAmount,
		&order.Currency,
		&order.DeliveryAddr,
		&order.Notes,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	// Fetch items
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, order_id, item_name, quantity, unit_price, total_price
		FROM order_items WHERE order_id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.OrderItem
		if err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ItemName,
			&item.Quantity,
			&item.UnitPrice,
			&item.TotalPrice,
		); err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}

	return &order, nil
}

func (r *PostgresOrderRepo) UpdateStatus(
	ctx context.Context,
	id uuid.UUID,
	expectedCurrent,
	newStatus model.OrderStatus,
) error {

	res, err := r.db.ExecContext(ctx, `
		UPDATE orders
		SET status = $1, updated_at = NOW()
		WHERE id = $2 AND status = $3
	`,
		newStatus,
		id,
		expectedCurrent,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrStaleUpdate
	}

	return nil
}

func (r *PostgresOrderRepo) ListByUser(
	ctx context.Context,
	userID uuid.UUID,
) ([]model.Order, error) {

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, status, total_amount, currency, delivery_addr, notes, created_at, updated_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order

	for rows.Next() {
		var o model.Order

		err := rows.Scan(
			&o.ID,
			&o.UserID,
			&o.Status,
			&o.TotalAmount,
			&o.Currency,
			&o.DeliveryAddr,
			&o.Notes,
			&o.CreatedAt,
			&o.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		orders = append(orders, o)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}
