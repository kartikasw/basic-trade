package common

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	ErrUniqueViolation = "23505"
)

var ErrRecordNotFound = pgx.ErrNoRows

type Error int

const (
	ErrCredentiials Error = iota
)

func ErrorCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}

	return err.Error()
}

func ErrorValidation(e error) error {
	errors := e.(validator.ValidationErrors)
	var err error

	for _, e := range errors {
		switch e.Tag() {
		case "required":
			err = fmt.Errorf("%s is required.", e.Field())
		case "min":
			if e.Kind() == reflect.Int {
				err = fmt.Errorf("%s must be at least %s.", e.Field(), e.Param())
			} else {
				err = fmt.Errorf("%s must be at least %s characters.", e.Field(), e.Param())
			}
		default:
			err = fmt.Errorf("%s is invalid.", e.Field())
		}
	}

	return err
}
