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
	err := e

	if errors, ok := e.(validator.ValidationErrors); ok {
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
			case "max":
				if e.Kind() == reflect.Int {
					err = fmt.Errorf("%s must be less than %s.", e.Field(), e.Param())
				} else {
					err = fmt.Errorf("%s must be less than %s characters.", e.Field(), e.Param())
				}
			case "validImage":
				err = fmt.Errorf("Image must be a jpg, jpeg, png, or svg, and no larger than 5 MB.")
			default:
				err = fmt.Errorf("%s is invalid.", e.Field())
			}
		}
	}

	return err
}
