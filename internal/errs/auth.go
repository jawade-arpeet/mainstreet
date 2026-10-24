package errs

import "errors"

const (
	CodeInvalidEmailOrPassword = "INVALID_EMAIL_OR_PASSWORD"
	CodeAccAlreadyExists       = "ACCOUNT_ALREADY_EXISTS"
)

var (
	ErrAccDoesNotExists = errors.New("account does not exists")
	ErrAccAlreadyExists = errors.New("account already exists")
	ErrInvalidPassword  = errors.New("invalid password")
)
