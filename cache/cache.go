// Package cache implements generic inline caches.
package cache

import (
	"fmt"

	"github.com/mrogaski/go-inline/internal/pkg/dictionary"
)

// BackingStoreAccessor is a generic interface defining the access methods for the backing store that the
// caching layer will be applied to.
type BackingStoreAccessor[K comparable, V any] interface {
	Get(K) (V, error)
	Set(K, V) error
}

// LRUCache is a cache implementation with a least-recently-used replacement policy.
type LRUCache[K comparable, V any] struct {
	pool    *dictionary.LinkedHashMap[K, V]
	maxSize int
	store   BackingStoreAccessor[K, V]
}

// NewLRUCache constructs a new LRUCache object with the maximum size set to maxSize.
func NewLRUCache[K comparable, V any](store BackingStoreAccessor[K, V], maxSize int) *LRUCache[K, V] {
	return &LRUCache[K, V]{pool: dictionary.NewLinkedHashMap[K, V](), maxSize: maxSize, store: store}
}

// Get returns either a cached entry for the key, if one is currently in the cached entry set, or an entry retrieved
// from the backing store.  If an error occurs on retrieval from the backing store, a wrapped error is returned.
func (c *LRUCache[K, V]) Get(key K) (V, error) {
	var err error

	result, ok := c.pool.Get(key)
	if !ok {
		result, err = c.store.Get(key)
		if err != nil {
			return result, fmt.Errorf("backing store error: %w", err)
		}
	}

	for c.pool.Size() >= c.maxSize {
		c.pool.EvictLeastRecent()
	}

	c.pool.Add(key, result)

	return result, nil
}
