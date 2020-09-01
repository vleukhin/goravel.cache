package store

import (
	"context"
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gudron/goravel.cache/errs"
	"github.com/gudron/goravel.cache/helpers"
)

// memcacheStore is a collection of method to access to store
type memcacheStore struct {
	mc     *memcache.Client
	prefix string
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

func (s *memcacheStore) Many(keys []string) (map[string][]byte, error) {
	prefixedKeys := make([]string, len(keys))
	for i:=0; i < len(keys); i++ {
		prefixedKeys[i]= s.prefix + keys[i]
	}

	items, err := s.mc.GetMulti(prefixedKeys)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, errs.ErrCacheMiss
	}

	result := make(map[string][]byte, 0)
	for key, item := range items {
		result[key] = item.Value
	}

	return result, nil
}

func (s *memcacheStore) Put(key string, value []byte, seconds int) error {
	item := &memcache.Item{
		Key:        key,
		Value:      value,
		Expiration: int32(s.calculateExpiration(seconds)),
		Flags:      0,
	}

	item.Key = s.prefix + item.Key

	err := s.mc.Set(item)
	if err != nil {
		return err
	}

	return nil
}

func (s *memcacheStore) PutMany(values map[string][]byte, seconds int) error {
	expiration := int32(s.calculateExpiration(seconds))

	for key, value := range values {
		item := &memcache.Item{
			Key:        key,
			Value:      value,
			Expiration: expiration,
			Flags:      0,
		}

		item.Key = s.prefix + item.Key

		err := s.mc.Set(item)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *memcacheStore) Add(key string, value []byte, seconds int) error {
	item := &memcache.Item{
		Key:        key,
		Value:      value,
		Expiration: int32(s.calculateExpiration(seconds)),
		Flags:      0,
	}

	err := s.mc.Add(item)
	if err != nil {
		return err
	}

	return nil
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

func (s *memcacheStore) Forever(key string, value []byte) (bool, error) {
	err := s.Put(key, value, 0)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *memcacheStore) GetPrefix() string {
	return s.prefix
}

func (s *memcacheStore) setPrefix(prefix string) {
	s.prefix = prefix + ":"
}

func (s *memcacheStore) Forget(key string) error {
	return s.mc.Delete(key)
}

func (s *memcacheStore) Flush() error {
	return s.mc.FlushAll()
}

func (s *memcacheStore) calculateExpiration(seconds int) int64 {
	return s.toTimestamp(seconds)
}

func (s *memcacheStore) toTimestamp(seconds int) int64 {
	if seconds > 0 {
		return helpers.AvailableAt(time.Duration(seconds) * time.Second)
	}

	return 0
}

// NewMemcacheStore on memcached
func NewMemcacheStore(ctx context.Context, cfg CacheStoreConfig) (*memcacheStore, error) {
	servers := make([]string, len(cfg.Hosts))
	for _, hostParam := range cfg.Hosts {
		servers = append(servers, fmt.Sprintf(
			"%s:%d",
			hostParam.Host, hostParam.Port))
	}

	mc := memcache.New(servers...)

	err := mc.Ping()
	if err != nil {
		return nil, err
	}

	mc.MaxIdleConns = cfg.MaxIdleConnections
	mc.Timeout = time.Duration(cfg.ReadWriteTimeOut) * time.Millisecond
	store := &memcacheStore{
		mc: mc,
	}

	store.setPrefix(cfg.Prefix)

	return store, nil
}
