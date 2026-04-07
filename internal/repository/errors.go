package repository

import (
	"github.com/pkg/errors"
)

var (
	errNestedTx = errors.New("nested transactions are not allowed")
)
