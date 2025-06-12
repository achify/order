package delivery

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Repository defines persistence behavior for deliveries and locations.
type Repository interface {
	Create(ctx context.Context, d *Delivery) error
	GetByID(ctx context.Context, id string) (*Delivery, error)
	List(ctx context.Context) ([]Delivery, error)
	Update(ctx context.Context, d *Delivery) error
	Delete(ctx context.Context, id string) error

	UpsertLocation(ctx context.Context, loc *Location) error
	Locations(ctx context.Context, provider string) ([]Location, error)
}

const createDeliveryQuery = `INSERT INTO deliveries
    (id, created_at, updated_at, provider, tracking_code, status)
    VALUES (:id, :created_at, :updated_at, :provider, :tracking_code, :status)`

const getDeliveryQuery = `SELECT * FROM deliveries WHERE id=$1`
const listDeliveryQuery = `SELECT * FROM deliveries`
const updateDeliveryQuery = `UPDATE deliveries SET
    updated_at=:updated_at,
    status=:status
    WHERE id=:id`
const deleteDeliveryQuery = `DELETE FROM deliveries WHERE id=$1`

const upsertLocationQuery = `INSERT INTO locations (id, provider, data, updated_at)
    VALUES (:id, :provider, :data, :updated_at)
    ON CONFLICT (id) DO UPDATE SET provider=EXCLUDED.provider, data=EXCLUDED.data, updated_at=EXCLUDED.updated_at`
const listLocationsQuery = `SELECT * FROM locations WHERE provider=$1`

// PostgresRepository implements Repository using PostgreSQL

type PostgresRepository struct {
	DB *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{DB: db}
}

func (r *PostgresRepository) Create(ctx context.Context, d *Delivery) error {
	_, err := r.DB.NamedExecContext(ctx, createDeliveryQuery, d)
	return err
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*Delivery, error) {
	var d Delivery
	if err := r.DB.GetContext(ctx, &d, getDeliveryQuery, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &d, nil
}

func (r *PostgresRepository) List(ctx context.Context) ([]Delivery, error) {
	var list []Delivery
	if err := r.DB.SelectContext(ctx, &list, listDeliveryQuery); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *PostgresRepository) Update(ctx context.Context, d *Delivery) error {
	_, err := r.DB.NamedExecContext(ctx, updateDeliveryQuery, d)
	return err
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	_, err := r.DB.ExecContext(ctx, deleteDeliveryQuery, id)
	return err
}

func (r *PostgresRepository) UpsertLocation(ctx context.Context, loc *Location) error {
	_, err := r.DB.NamedExecContext(ctx, upsertLocationQuery, loc)
	return err
}

func (r *PostgresRepository) Locations(ctx context.Context, provider string) ([]Location, error) {
	var l []Location
	if err := r.DB.SelectContext(ctx, &l, listLocationsQuery, provider); err != nil {
		return nil, err
	}
	return l, nil
}
