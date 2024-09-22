package httpServer

import (
	"github.com/pkg/errors"
)

var (
	errInternal = errors.New("internal error, please, try again later")
	errNotFound = errors.New("not found")
)
