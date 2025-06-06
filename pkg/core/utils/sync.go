package utils

import (
	"sync"
)

/**
 * MutexValue is a thread-safe wrapper around a value of type T.
 */
type MutexValue[T any] struct {
	v  T
	mu sync.Mutex
}

func (mv *MutexValue[T]) Get() T {
	mv.mu.Lock()
	defer mv.mu.Unlock()
	return mv.v
}

func (mv *MutexValue[T]) Set(value T) {
	mv.mu.Lock()
	defer mv.mu.Unlock()
	mv.v = value
}

func NewMutexValue[T any](value T) *MutexValue[T] {
	return &MutexValue[T]{v: value}
}
