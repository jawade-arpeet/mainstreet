package repository

import (
	"context"
	"mainstreet/internal/client"
	"mainstreet/internal/dao"

	"github.com/jackc/pgx/v5"
)

type AuthRepository struct {
	pg *client.PostgresClient
}

func newAuthRepository(pg *client.PostgresClient) *AuthRepository {
	return &AuthRepository{pg: pg}
}

func (r *AuthRepository) CreateAccount(
	ctx context.Context,
	email string,
	passwordHash string,
) (int64, error) {
	const query = `
		INSERT INTO accounts (
			email,
			password_hash
		)
		VALUES (
			@email,
			@password_hash
		);
	`
	args := pgx.NamedArgs{
		"email":         email,
		"password_hash": passwordHash,
	}

	return r.pg.Exec(ctx, query, args)
}

func (r *AuthRepository) GetAccountByEmail(
	ctx context.Context,
	email string,
) (dao.Account, error) {
	const query = `
		SELECT
			id,
			email,
			password_hash,
			is_active,
			is_email_verified
		FROM
			accounts
		WHERE
			email = @email
	`
	args := pgx.NamedArgs{"email": email}

	return r.pg.QueryOne[dao.Account](ctx, query, args)
}
