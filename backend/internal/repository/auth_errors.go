package repository

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"

	"github.com/Project-TETI/AgriSense/backend/internal/domain"
	"github.com/jackc/pgx/v5/pgconn"
)

func authDBError(err error) error {
	if err == nil || errors.Is(err, domain.ErrAuthUnavailable) {
		return err
	}
	var connectionError *pgconn.ConnectError
	var pgError *pgconn.PgError
	unavailable := errors.Is(err, context.DeadlineExceeded) ||
		errors.Is(err, driver.ErrBadConn) || errors.As(err, &connectionError)
	if errors.As(err, &pgError) {
		unavailable = unavailable || strings.HasPrefix(pgError.Code, "08") ||
			pgError.Code == "57P01" || pgError.Code == "57P02" || pgError.Code == "57P03"
	}
	if unavailable {
		return fmt.Errorf("%w: %w", domain.ErrAuthUnavailable, err)
	}
	return err
}
