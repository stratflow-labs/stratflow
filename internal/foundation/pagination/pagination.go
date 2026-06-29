package pagination

import (
	"errors"
	"math"
)

var ErrOffsetOverflow = errors.New("pagination offset overflow")

type Config struct {
	DefaultPage     int
	DefaultPageSize int
	MaxPageSize     int
}

type Value struct {
	Page     int
	PageSize int
	Offset   int
}

func Resolve(page, pageSize int, cfg Config) (Value, error) {
	defaultPage := cfg.DefaultPage
	if defaultPage <= 0 {
		defaultPage = 1
	}

	defaultPageSize := cfg.DefaultPageSize
	if defaultPageSize <= 0 {
		defaultPageSize = 20
	}

	maxPageSize := cfg.MaxPageSize
	if maxPageSize <= 0 {
		maxPageSize = defaultPageSize
	}

	resolvedPage := page
	if resolvedPage <= 0 {
		resolvedPage = defaultPage
	}

	resolvedPageSize := pageSize
	if resolvedPageSize <= 0 {
		resolvedPageSize = defaultPageSize
	}
	if resolvedPageSize > maxPageSize {
		resolvedPageSize = maxPageSize
	}

	offset := int64(resolvedPage-1) * int64(resolvedPageSize)
	if offset > math.MaxInt32 {
		return Value{}, ErrOffsetOverflow
	}

	return Value{
		Page:     resolvedPage,
		PageSize: resolvedPageSize,
		Offset:   int(offset),
	}, nil
}
