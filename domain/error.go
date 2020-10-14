package domain

import "errors"

var (
	ErrTransactionAlreadyExist = errors.New("transaction ID already exist")
	ErrTransactionEmptyID      = errors.New("transaction must have ID")
	ErrNotFound                = errors.New("resource not found")
)
