package basket

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Repository defines persistence for baskets and their items.
type Repository interface {
	Create(ctx context.Context, b *Basket) error
	GetByID(ctx context.Context, id string) (*Basket, error)
	List(ctx context.Context) ([]Basket, error)
	Update(ctx context.Context, b *Basket) error
	Delete(ctx context.Context, id string) error

	AddItem(ctx context.Context, it *Item) error
	UpdateItem(ctx context.Context, it *Item) error
	DeleteItem(ctx context.Context, basketID, itemID string) error
	ListItems(ctx context.Context, basketID string) ([]Item, error)
}

const createQuery = `INSERT INTO baskets
    (id, created_at, updated_at, account_id, total_price)
    VALUES (:id, :created_at, :updated_at, :account_id, :total_price)`
const getQuery = `SELECT * FROM baskets WHERE id=$1`
const listQuery = `SELECT * FROM baskets`
const updateQuery = `UPDATE baskets SET
    updated_at=:updated_at,
    account_id=:account_id,
    total_price=:total_price
    WHERE id=:id`
const deleteQuery = `DELETE FROM baskets WHERE id=$1`

const addItemQuery = `INSERT INTO basket_items
    (basket_id, item_id, quantity, price_per_item)
    VALUES ($1,$2,$3,$4)`
const updateItemQuery = `UPDATE basket_items SET quantity=$3, price_per_item=$4 WHERE basket_id=$1 AND item_id=$2`
const deleteItemQuery = `DELETE FROM basket_items WHERE basket_id=$1 AND item_id=$2`
const listItemsQuery = `SELECT * FROM basket_items WHERE basket_id=$1`

// PostgresRepository implements Repository using PostgreSQL.
type PostgresRepository struct{ DB *sqlx.DB }

func NewPostgresRepository(db *sqlx.DB) *PostgresRepository { return &PostgresRepository{DB: db} }

func (r *PostgresRepository) Create(ctx context.Context, b *Basket) error {
	_, err := r.DB.NamedExecContext(ctx, createQuery, b)
	return err
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*Basket, error) {
	var b Basket
	if err := r.DB.GetContext(ctx, &b, getQuery, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &b, nil
}

func (r *PostgresRepository) List(ctx context.Context) ([]Basket, error) {
	var list []Basket
	if err := r.DB.SelectContext(ctx, &list, listQuery); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *PostgresRepository) Update(ctx context.Context, b *Basket) error {
	_, err := r.DB.NamedExecContext(ctx, updateQuery, b)
	return err
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	_, err := r.DB.ExecContext(ctx, deleteQuery, id)
	return err
}

func (r *PostgresRepository) AddItem(ctx context.Context, it *Item) error {
	_, err := r.DB.ExecContext(ctx, addItemQuery, it.BasketID, it.ItemID, it.Quantity, it.PricePerItem)
	return err
}

func (r *PostgresRepository) UpdateItem(ctx context.Context, it *Item) error {
	_, err := r.DB.ExecContext(ctx, updateItemQuery, it.BasketID, it.ItemID, it.Quantity, it.PricePerItem)
	return err
}

func (r *PostgresRepository) DeleteItem(ctx context.Context, basketID, itemID string) error {
	_, err := r.DB.ExecContext(ctx, deleteItemQuery, basketID, itemID)
	return err
}

func (r *PostgresRepository) ListItems(ctx context.Context, basketID string) ([]Item, error) {
	var list []Item
	if err := r.DB.SelectContext(ctx, &list, listItemsQuery, basketID); err != nil {
		return nil, err
	}
	return list, nil
}
