package pool_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mrogaski/go-inline/internal/pkg/pool"
)

func TestNewLinkedMap(t *testing.T) {
	t.Parallel()

	got := pool.NewLinkedHashMap[string, string]()

	assert.NotEmpty(t, got)
}

func TestLinkedMap_empty(t *testing.T) {
	t.Parallel()

	m := pool.NewLinkedHashMap[string, string]()

	got, ok := m.Get("A")

	assert.Empty(t, got)
	assert.False(t, ok)
	assert.Equal(t, 0, m.Size())
}

func TestLinkedMap_single_element(t *testing.T) {
	t.Parallel()

	m := pool.NewLinkedHashMap[string, string]()

	m.Add("A", "alpha")

	got, ok := m.Get("A")

	assert.Equal(t, got, "alpha")
	assert.True(t, ok)
	assert.Equal(t, 1, m.Size())

	got, ok = m.Get("Z")

	assert.Empty(t, got)
	assert.False(t, ok)
}

func TestLinkedMap_multiple_elements(t *testing.T) {
	t.Parallel()

	m := pool.NewLinkedHashMap[string, string]()

	m.Add("A", "alpha")
	m.Add("B", "bravo")
	m.Add("C", "charlie")

	got, ok := m.Get("B")

	assert.Equal(t, got, "bravo")
	assert.True(t, ok)
	assert.Equal(t, 3, m.Size())

	got, ok = m.Get("Z")

	assert.Empty(t, got)
	assert.False(t, ok)
}

func TestLinkedMap_overwrite_element(t *testing.T) {
	t.Parallel()

	m := pool.NewLinkedHashMap[string, string]()

	m.Add("A", "alpha")
	m.Add("B", "baker")
	m.Add("C", "charlie")
	m.Add("B", "bravo")

	got, ok := m.Get("B")

	assert.Equal(t, got, "bravo")
	assert.True(t, ok)
	assert.Equal(t, 3, m.Size())
}

func TestLinkedMap_EvictLeastRecent(t *testing.T) {
	t.Parallel()

	m := pool.NewLinkedHashMap[string, string]()

	m.Add("A", "alpha")
	m.Add("B", "bravo")
	m.Add("C", "charlie")

	_ = m.EvictLeastRecent()

	assert.Equal(t, 2, m.Size())

	got, ok := m.Get("A")

	assert.Empty(t, got)
	assert.False(t, ok)

	got, ok = m.Get("B")

	assert.Equal(t, got, "bravo")
	assert.True(t, ok)
}
