package domain

import (
	"errors"
	"fmt"
)

var (
	ErrValidation = errors.New("validation error")

	ErrPageLimit = fmt.Errorf("invalid page limit [max=%v]", MaxPageLimit)

	ErrBookName     = errors.New("invalid name")
	ErrBookNotFound = errors.New("book not found")
)
