package database

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrNotFound       = errors.New("not found")
	ErrWrongInputData = errors.New("wrong input data")
	ErrAlreadyExists  = errors.New("already exists")
)

// Error extracting database error
func Error(err error) error {
	if err == nil {
		return nil
	}
	pgErr := new(pgconn.PgError)
	isPgErr := errors.As(err, &pgErr)
	if isPgErr {
		switch pgErr.Code {
		case "23502":
			return ErrWrongInputData
		case "23505":
			return ErrAlreadyExists
		case "42P07":
			return ErrAlreadyExists
		}
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	return err
}
