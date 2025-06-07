package payment

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Repository defines persistence layer for payments.
type Repository interface {
	Create(ctx context.Context, p *Payment) error
	GetByID(ctx context.Context, id string) (*Payment, error)
}

const createQuery = `INSERT INTO payments
    (id, created_at, updated_at, order_id, amount, status)
    VALUES (:id, :created_at, :updated_at, :order_id, :amount, :status)`

const getQuery = `SELECT * FROM payments WHERE id=$1`

type PostgresRepository struct{ DB *sqlx.DB }

func NewPostgresRepository(db *sqlx.DB) *PostgresRepository { return &PostgresRepository{DB: db} }

func (r *PostgresRepository) Create(ctx context.Context, p *Payment) error {
	_, err := r.DB.NamedExecContext(ctx, createQuery, p)
	return err
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*Payment, error) {
	var p Payment
	if err := r.DB.GetContext(ctx, &p, getQuery, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}
