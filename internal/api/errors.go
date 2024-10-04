package api

import (
	"github.com/pkg/errors"
)

var (
	ErrInternal   = errors.New("internal error, please, try again later")
	ErrNotFound   = errors.New("not found")
	ErrWrongInput = errors.New("fill all required fields")
)
