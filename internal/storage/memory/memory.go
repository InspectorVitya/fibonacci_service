package memory

import (
	"context"
	"sync"

	"github.com/inspectorvitya/fibonacci_service/internal/storage"
)

type StorageInMemory struct {
	mu       sync.RWMutex
	fibSlice []string
}

func New() storage.Storage {
	return &StorageInMemory{
		fibSlice: make([]string, 100),
	}
}

func (s *StorageInMemory) GetSliceFib(_ context.Context, x, y int) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.fibSlice[x-1 : y], nil
}

func (s *StorageInMemory) SetSliceFib(_ context.Context, fibSlice []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.fibSlice = fibSlice
	return nil
}
