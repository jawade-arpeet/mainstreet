package client

import (
	"context"
	"errors"
	"fmt"
	"mainstreet/internal/config"
	"mainstreet/internal/errs"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresClient struct {
	pool *pgxpool.Pool
}

var errPgOperationFailed = errors.New("postgres operation failed")

func newPostgresClient(ctx context.Context) (*PostgresClient, error) {
	cfg := config.GetPostgresConfig()

	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return &PostgresClient{pool: pool}, nil
}

func (c *PostgresClient) Close() {
	c.pool.Close()
}

func (c *PostgresClient) wrapErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("%w: %w", errs.ErrPgNoRows, err)
	}

	pgErr, ok := errors.AsType[*pgconn.PgError](err)
	if !ok {
		return fmt.Errorf("%w: %w", errPgOperationFailed, err)
	}

	var sentinel error
	switch pgErr.Code {
	case pgerrcode.UniqueViolation:
		sentinel = errs.ErrPgUniqueViolation
	case pgerrcode.ForeignKeyViolation:
		sentinel = errs.ErrPgForeignKeyViolation
	case pgerrcode.NotNullViolation:
		sentinel = errs.ErrPgNotNullViolation
	case pgerrcode.CheckViolation:
		sentinel = errs.ErrPgCheckViolation
	default:
		return fmt.Errorf("%w: %w", errPgOperationFailed, err)
	}

	return fmt.Errorf("%w: %w", sentinel, err)
}

func (c *PostgresClient) Exec(
	ctx context.Context,
	query string,
	args pgx.NamedArgs,
) (int64, error) {
	tag, err := c.pool.Exec(ctx, query, args)
	if err != nil {
		return 0, c.wrapErr(err)
	}

	return tag.RowsAffected(), nil
}

func (c *PostgresClient) QueryOne[T any](
	ctx context.Context,
	query string,
	args pgx.NamedArgs,
) (T, error) {
	var zero T

	rows, err := c.pool.Query(ctx, query, args)
	if err != nil {
		return zero, c.wrapErr(err)
	}

	defer rows.Close()

	var result T
	if err := pgxscan.ScanOne(&result, rows); err != nil {
		return zero, c.wrapErr(err)
	}

	return result, nil
}

func (c *PostgresClient) QueryAll[T any](
	ctx context.Context,
	query string,
	args pgx.NamedArgs,
) ([]T, error) {
	rows, err := c.pool.Query(ctx, query, args)
	if err != nil {
		return nil, c.wrapErr(err)
	}

	defer rows.Close()

	var result []T
	if err := pgxscan.ScanAll(&result, rows); err != nil {
		return nil, c.wrapErr(err)
	}

	return result, nil
}
