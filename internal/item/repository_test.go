package item

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestRepositoryCreate(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	xdb := sqlx.NewDb(db, "sqlmock")
	repo := NewPostgresRepository(xdb)

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO items")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "name", int64(1), "c").
		WillReturnResult(sqlmock.NewResult(1, 1))

	it := &Item{ID: "id", Name: "name", Price: 1, CategoryID: "c"}
	err := repo.Create(context.Background(), it)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expect: %v", err)
	}
}
