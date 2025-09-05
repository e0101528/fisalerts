package repositories

import (
	"context"
)

type IFluxStore interface {
	RunQuery(ctx context.Context) error
	RunQueryWithParams(ctx context.Context) error
}

type ICheckState interface {
	SetActive(ctx context.Context) error
	IsActive(ctx context.Context) error
	SetClearable()
	Clear()
}
