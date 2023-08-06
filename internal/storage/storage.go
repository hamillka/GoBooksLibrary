package storage

import "errors"

var (
	ErrBookNotFound = errors.New("book not found")
	ErrBookExists   = errors.New("book exists")
)
