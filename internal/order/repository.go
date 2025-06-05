package order

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Repository defines persistence behavior for orders

type Repository interface {
	Create(ctx context.Context, o *Order) error
	GetByID(ctx context.Context, id string) (*Order, error)
	// List returns orders optionally filtered by delivery id
	List(ctx context.Context, deliveryID string) ([]Order, error)
	Update(ctx context.Context, o *Order) error
	Delete(ctx context.Context, id string) error
}

// PostgresRepository implements Repository using PostgreSQL

const createQuery = `INSERT INTO orders
    (id, created_at, updated_at, receiver_id, account_id, seller_id, delivery_id, basket_id)
    VALUES (:id, :created_at, :updated_at, :receiver_id, :account_id, :seller_id, :delivery_id, :basket_id)`

const getQuery = `SELECT * FROM orders WHERE id=$1`
const listQuery = `SELECT * FROM orders`
const listByDeliveryQuery = `SELECT * FROM orders WHERE delivery_id=$1`
const updateQuery = `UPDATE orders SET
    updated_at=:updated_at,
    receiver_id=:receiver_id,
    account_id=:account_id,
    seller_id=:seller_id,
    delivery_id=:delivery_id,
    basket_id=:basket_id
    WHERE id=:id`
const deleteQuery = `DELETE FROM orders WHERE id=$1`

type PostgresRepository struct {
	DB *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) *PostgresRepository { return &PostgresRepository{DB: db} }

func (r *PostgresRepository) Create(ctx context.Context, o *Order) error {
	_, err := r.DB.NamedExecContext(ctx, createQuery, o)
	return err
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*Order, error) {
	var o Order
	if err := r.DB.GetContext(ctx, &o, getQuery, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &o, nil
}

func (r *PostgresRepository) List(ctx context.Context, deliveryID string) ([]Order, error) {
	var list []Order
	query := listQuery
	args := []interface{}{}
	if deliveryID != "" {
		query = listByDeliveryQuery
		args = append(args, deliveryID)
	}
	if err := r.DB.SelectContext(ctx, &list, query, args...); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *PostgresRepository) Update(ctx context.Context, o *Order) error {
	_, err := r.DB.NamedExecContext(ctx, updateQuery, o)
	return err
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	_, err := r.DB.ExecContext(ctx, deleteQuery, id)
	return err
}
