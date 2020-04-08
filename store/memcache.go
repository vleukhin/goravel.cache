package store

import (
	"context"
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gudron/goravel.cache/errs"
)

// memcacheStore is a collection of method to access to store
type memcacheStore struct {
	mc     *memcache.Client
	prefix string
}

func (s *memcacheStore) Forever(key string, value []byte) (bool, error) {
	item := &memcache.Item{
		Key:        key,
		Value:      value,
		Expiration: 0,
		Flags:      0,
	}

	err := s.put(item)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *memcacheStore) put(value *memcache.Item) error {
	value.Key = s.prefix + value.Key

	err := s.mc.Set(value)
	if err != nil {
		return err
	}

	return nil
}

func (s *memcacheStore) Get(key string) ([]byte, error) {
	prefixedKey := s.prefix + key

	item, err := s.mc.Get(prefixedKey)
	if err != nil && err == memcache.ErrCacheMiss {
		return nil, errs.ErrCacheMiss
	}

	if item == nil {
		return nil, errs.ErrCacheMiss
	}

	return item.Value, nil
}

func (s *memcacheStore) Increment(key string, value uint64) (uint64, error) {
	prefixedKey := s.prefix + key

	newValue, err := s.mc.Increment(prefixedKey, value)
	if err != nil {
		return 0, err
	}

	return newValue, nil
}

func (s *memcacheStore) Decrement(key string, value uint64) (uint64, error) {
	prefixedKey := s.prefix + key

	newValue, err := s.mc.Decrement(prefixedKey, value)
	if err != nil {
		return 0, err
	}

	return newValue, nil
}

func (s *memcacheStore) GetPrefix() string {
	return s.prefix
}

func (s *memcacheStore) setPrefix(prefix string) {
	s.prefix = prefix + ":"
}

// NewMemcacheStore on memcached
func NewMemcacheStore(ctx context.Context, host string, port int, prefix string) (*memcacheStore, error) {
	mc := memcache.New(fmt.Sprintf(
		"%s:%d",
		host, port))

	store := &memcacheStore{
		mc: mc,
	}

	store.setPrefix(prefix)

	return store, nil
}
