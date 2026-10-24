package service

import (
	"context"
	"errors"
	"fmt"
	"mainstreet/internal/argon2id"
	"mainstreet/internal/errs"
	"mainstreet/internal/repository"
	"strings"
)

type AuthService struct {
	authRepo *repository.AuthRepository
}

func newAuthService(authRepo *repository.AuthRepository) *AuthService {
	return &AuthService{authRepo: authRepo}
}

func (s *AuthService) SignUp(
	ctx context.Context,
	email string,
	password string,
) error {
	email = strings.ToLower(email)
	email = strings.TrimSpace(email)
	password = strings.TrimSpace(password)

	passwordHash, err := argon2id.Hash(password)
	if err != nil {
		return errs.Internal(fmt.Errorf("failed to hash password: %w", err))
	}

	rowsAffected, err := s.authRepo.CreateAccount(ctx, email, passwordHash)
	if err != nil {
		if errors.Is(err, errs.ErrPgUniqueViolation) {
			return errs.Conflict(
				errs.CodeAccAlreadyExists,
				"account already exists",
				errs.ErrAccAlreadyExists,
			)
		}

		return errs.Internal(fmt.Errorf("failed to create account: %w", err))
	}

	if rowsAffected != 1 {
		return errs.Internal(fmt.Errorf("failed to create account: %w", err))
	}

	return nil
}

func (s *AuthService) SignIn(
	ctx context.Context,
	email string,
	password string,
) (string, error) {
	email = strings.ToLower(email)
	email = strings.TrimSpace(email)
	password = strings.TrimSpace(password)

	acc, err := s.authRepo.GetAccountByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, errs.ErrPgNoRows) {
			return "", errs.Unauthorized(
				errs.CodeInvalidEmailOrPassword,
				"invalid email or password",
				errs.ErrAccDoesNotExists,
			)
		}

		return "", errs.Internal(fmt.Errorf("failed to retrieve account: %w", err))
	}

	ok, err := argon2id.Compare(password, acc.PasswordHash)
	if err != nil {
		if errors.Is(err, argon2id.ErrInvalidHash) {
			return "", errs.Unauthorized(
				errs.CodeInvalidEmailOrPassword,
				"invalid email or password",
				errs.ErrAccDoesNotExists,
			)
		}

		return "", errs.Internal(fmt.Errorf("failed to retrieve account: %w", err))
	}

	if !ok {
		return "", errs.Unauthorized(
			errs.CodeInvalidEmailOrPassword,
			"invalid email or password",
			errs.ErrAccDoesNotExists,
		)
	}

	return acc.ID.String(), nil
}
