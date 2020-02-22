package store

import (
	"context"

	"gitlab.com/nativerent/adunit/pkg/common/cache/errs"
)

// inMemoryStore is a collection of method to for run tests
type inMemoryStore struct {
	prefix string
	data   map[string][]byte
}

func (s *inMemoryStore) Forever(key string, value []byte) (bool, error) {
	err := s.put(key, value)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *inMemoryStore) put(key string, value []byte) error {
	prefixedKey := s.prefix + key

	s.data[prefixedKey] = value

	return nil
}

func (s *inMemoryStore) Get(key string) ([]byte, error) {
	prefixedKey := s.prefix + key

	value, ok := s.data[prefixedKey]
	if !ok {
		return nil, errs.ErrCacheMiss
	}

	return value, nil
}

func (s *inMemoryStore) Increment(key string, value uint64) (uint64, error) {

	return 1, nil
}

func (s *inMemoryStore) Decrement(key string, value uint64) (uint64, error) {

	return 1, nil
}

func (s *inMemoryStore) GetPrefix() string {
	return s.prefix
}

// NewInMemoryStore on memcached
func NewInMemoryStore(ctx context.Context, prefix string) (*inMemoryStore, error) {
	store := &inMemoryStore{
		data:   make(map[string][]byte),
		prefix: prefix + ":",
	}

	return store, nil
}
