package repository

import (
	"context"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/Project-TETI/AgriSense/backend/internal/domain"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestAuthDBError(t *testing.T) {
	for _, cause := range []error{
		context.DeadlineExceeded, driver.ErrBadConn,
		&pgconn.PgError{Code: "08006"}, &pgconn.PgError{Code: "57P01"},
	} {
		got := authDBError(cause)
		if !errors.Is(got, domain.ErrAuthUnavailable) || !errors.Is(got, cause) {
			t.Fatalf("expected storage-unavailable classification preserving cause %T", cause)
		}
	}
	for _, cause := range []error{nil, context.Canceled, errors.New("unexpected failure"), domain.ErrInvalidRefreshToken} {
		if authDBError(cause) != cause {
			t.Fatalf("unexpected classification for %T", cause)
		}
	}
}
