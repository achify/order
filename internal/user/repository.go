package user

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Repository defines persistence behavior for users

type Repository interface {
	Create(ctx context.Context, u *User) error
	GetByUsername(ctx context.Context, username string) (*User, error)
	Update(ctx context.Context, u *User) error
}

// PostgresRepository implements Repository using PostgreSQL

const createQuery = `INSERT INTO users
    (id, created_at, updated_at, username, password, roles, active, seller_id, account_id)
    VALUES (:id, :created_at, :updated_at, :username, :password, :roles, :active, :seller_id, :account_id)`

const getQuery = `SELECT * FROM users WHERE username=$1 AND active=true`

const updateQuery = `UPDATE users SET
    updated_at=:updated_at,
    password=:password,
    roles=:roles,
    active=:active,
    seller_id=:seller_id,
    account_id=:account_id
    WHERE id=:id`

type PostgresRepository struct {
	DB *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) *PostgresRepository { return &PostgresRepository{DB: db} }

func (r *PostgresRepository) Create(ctx context.Context, u *User) error {
	_, err := r.DB.NamedExecContext(ctx, createQuery, u)
	return err
}

func (r *PostgresRepository) GetByUsername(ctx context.Context, username string) (*User, error) {
	var u User
	if err := r.DB.GetContext(ctx, &u, getQuery, username); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *PostgresRepository) Update(ctx context.Context, u *User) error {
	_, err := r.DB.NamedExecContext(ctx, updateQuery, u)
	return err
}
