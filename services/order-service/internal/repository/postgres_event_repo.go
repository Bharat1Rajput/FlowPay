package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type PostgresProcessedEventRepo struct {
	db *sql.DB
}

func NewProcessedEventRepo(db *sql.DB) *PostgresProcessedEventRepo {
	return &PostgresProcessedEventRepo{db: db}
}

func (r *PostgresProcessedEventRepo) Exists(ctx context.Context, eventID uuid.UUID) (bool, error) {
	var exists bool

	query := `SELECT EXISTS (SELECT 1 FROM processed_events WHERE event_id = $1)`
	err := r.db.QueryRowContext(ctx, query, eventID).Scan(&exists)

	return exists, err
}

func (r *PostgresProcessedEventRepo) Save(ctx context.Context, eventID uuid.UUID) error {
	query := `INSERT INTO processed_events (event_id) VALUES ($1)`
	_, err := r.db.ExecContext(ctx, query, eventID)

	return err
}
