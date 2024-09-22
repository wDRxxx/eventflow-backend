package service

import (
	"github.com/pkg/errors"
)

var (
	ErrPricesForFree     = errors.New("prices are provided for free event")
	ErrNoPrices          = errors.New("prices aren't provided for non-free event")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrWrongCredentials  = errors.New("wrong credentials")
	ErrPermissionDenied  = errors.New("permission denied")
)
