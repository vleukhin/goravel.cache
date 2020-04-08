package cache

// Store is interface to cache store
type Store interface {
	Get(key string) ([]byte, error)
	Increment(key string, value uint64) (uint64, error)
	Decrement(key string, value uint64) (uint64, error)
	Forever(key string, value []byte) (bool, error)
	GetPrefix() string
}

// Service is interface to cache service
type Service interface {
	Get(keys ...string) ([]byte, error)
	Increment(key string, value uint64) (uint64, error)
	Decrement(key string, value uint64) (uint64, error)
	Tags(keys ...string) (*taggedCacheService, error)
}

// service ...
type cacheService struct {
	Store Store
}

// Get ...
func (s *cacheService) Get(keys ...string) ([]byte, error) {
	keyNames := make([]string, len(keys))

	itemKey, err := s.itemKey(keyNames[0])
	if err != nil {
		return nil, err
	}

	value, err := s.Store.Get(itemKey)
	if err != nil {
		return nil, err
	}

	return value, nil
}

// Increment ...
func (s *cacheService) Increment(key string, value uint64) (uint64, error) {
	return 0, nil
}

// Decrement ...
func (s *cacheService) Decrement(key string, value uint64) (uint64, error) {
	return 0, nil
}

func (s *cacheService) itemKey(key string) (string, error) {
	return key, nil
}

// Tags ...
func (s *cacheService) Tags(keys ...string) (*taggedCacheService, error) {
	taggedCache, err := NewTaggedCacheService(s.Store, keys...)
	if err != nil {
		return nil, err
	}

	return taggedCache, nil
}

// New instance of simply cache
func NewCacheService(store Store) (*cacheService, error) {

	return &cacheService{Store: store}, nil
}
