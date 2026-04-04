package repository

import (
	"context"
	"database/sql"

	"github.com/Bharat1Rajput/flowpay/services/payment-service/internal/model"
	"github.com/google/uuid"
)

type PostgresPaymentRepo struct {
	db *sql.DB
}

func NewPostgresPaymentRepo(db *sql.DB) *PostgresPaymentRepo {
	return &PostgresPaymentRepo{db: db}
}

func (r *PostgresPaymentRepo) CreatePayment(ctx context.Context, p *model.Payment) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO payments (id, order_id, idempotency_key, amount, currency, status)
		VALUES ($1, $2, $3, $4, $5, $6)
	`,
		p.ID,
		p.OrderID,
		p.IdempotencyKey,
		p.Amount,
		p.Currency,
		p.Status,
	)
	return err
}

func (r *PostgresPaymentRepo) UpdateStatus(
	ctx context.Context,
	id uuid.UUID,
	expectedCurrent,
	newStatus model.PaymentStatus,
	gatewayRef,
	failureReason string,
) error {

	res, err := r.db.ExecContext(ctx, `
		UPDATE payments
		SET status = $1, gateway_ref = $2, failure_reason = $3, updated_at = NOW()
		WHERE id = $4 AND status = $5
	`,
		newStatus,
		gatewayRef,
		failureReason,
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
func (r *PostgresPaymentRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Payment, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, order_id, idempotency_key, amount, currency, status, gateway_ref, failure_reason, created_at, updated_at
		FROM payments WHERE id = $1
	`, id)

	var p model.Payment

	err := row.Scan(
		&p.ID,
		&p.OrderID,
		&p.IdempotencyKey,
		&p.Amount,
		&p.Currency,
		&p.Status,
		&p.GatewayRef,
		&p.FailureReason,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}
func (r *PostgresPaymentRepo) GetByIdempotencyKey(ctx context.Context, key string) (*model.Payment, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, order_id, idempotency_key, amount, currency, status, gateway_ref, failure_reason, created_at, updated_at
		FROM payments WHERE idempotency_key = $1
	`, key)

	var p model.Payment

	err := row.Scan(
		&p.ID,
		&p.OrderID,
		&p.IdempotencyKey,
		&p.Amount,
		&p.Currency,
		&p.Status,
		&p.GatewayRef,
		&p.FailureReason,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}
