package item

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Repository defines persistence behavior for items.
type Repository interface {
	Create(ctx context.Context, i *Item) error
	GetByID(ctx context.Context, id string) (*Item, error)
	List(ctx context.Context) ([]Item, error)
	Update(ctx context.Context, i *Item) error
	Delete(ctx context.Context, id string) error
}

const createQuery = `INSERT INTO items
    (id, created_at, updated_at, name, price, category_id)
    VALUES (:id, :created_at, :updated_at, :name, :price, :category_id)`
const getQuery = `SELECT * FROM items WHERE id=$1`
const listQuery = `SELECT * FROM items`
const updateQuery = `UPDATE items SET
    updated_at=:updated_at,
    name=:name,
    price=:price,
    category_id=:category_id
    WHERE id=:id`
const deleteQuery = `DELETE FROM items WHERE id=$1`

// PostgresRepository implements Repository using PostgreSQL.
type PostgresRepository struct {
	DB *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) *PostgresRepository { return &PostgresRepository{DB: db} }

func (r *PostgresRepository) Create(ctx context.Context, i *Item) error {
	_, err := r.DB.NamedExecContext(ctx, createQuery, i)
	return err
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*Item, error) {
	var it Item
	if err := r.DB.GetContext(ctx, &it, getQuery, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &it, nil
}

func (r *PostgresRepository) List(ctx context.Context) ([]Item, error) {
	var list []Item
	if err := r.DB.SelectContext(ctx, &list, listQuery); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *PostgresRepository) Update(ctx context.Context, i *Item) error {
	_, err := r.DB.NamedExecContext(ctx, updateQuery, i)
	return err
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	_, err := r.DB.ExecContext(ctx, deleteQuery, id)
	return err
}
