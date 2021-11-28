package storage

import "context"

type Storage interface {
	GetSliceFib(ctx context.Context, x, y int) ([]string, error)
	SetSliceFib(ctx context.Context, fibSlice []string) error
}
