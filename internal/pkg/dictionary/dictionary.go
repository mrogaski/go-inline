// Package dictionary is a cache entry pool implemented as a linked hash map.
package dictionary

import (
	"container/list"
	"fmt"
	"sync"
)

// LinkedHashMap is a hash map with the value elements organized as a double linked list.
type LinkedHashMap[K comparable, V any] struct {
	mutex     sync.Mutex
	entryMap  map[K]*list.Element
	entryList *list.List
}

// Entry holds the key/value pairs as the element content.
type Entry[K comparable, V any] struct {
	key   K
	value V
}

// NewLinkedHashMap constructs and initializes a new LinkedHashMap.
func NewLinkedHashMap[K comparable, V any]() *LinkedHashMap[K, V] {
	return &LinkedHashMap[K, V]{
		entryMap:  make(map[K]*list.Element),
		entryList: list.New(),
	}
}

// Size returns the number of entries contained in the LinkedHashMap.
func (m *LinkedHashMap[K, V]) Size() int {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.entryList.Len()
}

// Add inserts an entry into the LinkedHashMap.
func (m *LinkedHashMap[K, V]) Add(key K, value V) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	entry := Entry[K, V]{
		key:   key,
		value: value,
	}

	elem, ok := m.entryMap[key]
	if ok {
		elem.Value = entry
		m.entryList.MoveToBack(elem)
	} else {
		elem = m.entryList.PushBack(entry)
		m.entryMap[key] = elem
	}
}

// Get accesses a value by key and updates the position in the linked list.
func (m *LinkedHashMap[K, V]) Get(key K) (V, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var result V

	elem, ok := m.entryMap[key]
	if !ok {
		return result, false
	}

	entry, ok := elem.Value.(Entry[K, V])
	if !ok {
		return result, false
	}

	return entry.value, true
}

// EvictLeastRecent removes the entry which has been accessed least recently.
func (m *LinkedHashMap[K, V]) EvictLeastRecent() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	elem := m.entryList.Front()
	if elem == nil {
		return nil
	}

	entry, ok := elem.Value.(Entry[K, V])
	if !ok {
		return fmt.Errorf("invalid pool entry: %v", elem.Value)
	}

	m.entryList.Remove(elem)

	delete(m.entryMap, entry.key)

	return nil
}
