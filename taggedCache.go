package cache

import (
	"crypto/sha1"
	"fmt"
)

// taggedCacheService ...
type taggedCacheService struct {
	cacheService
	TagSet *TagSet
	Store  Store
}

func (tc *taggedCacheService) putManyForever() {

}

func (tc *taggedCacheService) putMany() {

}

func (tc *taggedCacheService) Get(keys ...string) ([]byte, error) {
	itemKey, err := tc.itemKey(keys[0])
	if err != nil {
		return nil, err
	}

	cacheData, err := tc.Store.Get(itemKey)
	if err != nil {
		return nil, err
	}

	return cacheData, nil
}

func (tc *taggedCacheService) Increment(key string, value uint64) (uint64, error) {
	newValue, err := tc.Store.Increment(key, value)
	if err != nil {
		return 0, err
	}

	return newValue, nil
}

func (tc *taggedCacheService) Decrement(key string, value uint64) (uint64, error) {
	newValue, err := tc.Store.Decrement(key, value)
	if err != nil {
		return 0, err
	}

	return newValue, nil
}

func (tc *taggedCacheService) flush() error {
	err := tc.TagSet.reset()
	if err != nil {
		return err
	}

	return nil
}

func (tc *taggedCacheService) itemKey(key string) (string, error) {
	k, err := tc.taggedItemKey(key)
	if err != nil {
		return "", err
	}

	return k, nil
}

func (tc *taggedCacheService) taggedItemKey(key string) (string, error) {
	namespace, err := tc.TagSet.getNamespace()
	if err != nil {
		return "", err
	}

	hash := sha1.Sum([]byte(namespace))
	hashString := fmt.Sprintf("%x", hash)
	result := hashString + ":" + key

	return result, nil

}

func (tc *taggedCacheService) getTags() *TagSet {
	return tc.TagSet
}

func (tc *taggedCacheService) Tags(keys ...string) (*taggedCacheService, error) {
	taggedCacheService, err := NewTaggedCacheService(tc.Store, keys...)
	if err != nil {
		return nil, err
	}

	return taggedCacheService, nil
}

// NewTaggedCacheService instance of tagged cache
func NewTaggedCacheService(store Store, names ...string) (*taggedCacheService, error) {
	tagSet := NewTagSet(store, names...)

	return &taggedCacheService{
		Store:  store,
		TagSet: tagSet,
	}, nil
}
