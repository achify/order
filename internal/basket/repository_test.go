package basket

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

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO baskets")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "a", int64(0)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	b := &Basket{ID: "id", AccountID: "a"}
	err := repo.Create(context.Background(), b)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expect: %v", err)
	}
}
