package cache_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/mrogaski/go-inline/cache"
)

type mockStore struct {
	mock.Mock
}

func (m *mockStore) Get(key string) (string, error) {
	args := m.Called(key)

	return args.String(0), args.Error(1)
}

func (m *mockStore) Set(key, value string) error {
	args := m.Called(key, value)

	return args.Error(0)
}

func TestNewLRUCache(t *testing.T) {
	t.Parallel()

	got := cache.NewLRUCache[string, string](&mockStore{}, 16)

	assert.NotNil(t, got)
}

func TestLRUCache_Get(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		key     string
		history []string
		store   *mockStore
		size    int
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "miss",
			key:  "A",
			store: func() *mockStore {
				m := &mockStore{}

				m.On("Get", "A").Return("automaton", nil).Once()

				return m
			}(),
			size:    4,
			want:    "automaton",
			wantErr: assert.NoError,
		},
		{
			name: "miss failure",
			key:  "A",
			store: func() *mockStore {
				m := &mockStore{}

				m.On("Get", "A").Return("", errors.New("failure")).Once()

				return m
			}(),
			size: 4,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "backing store error: failure")
			},
		},
		{
			name:    "hit",
			key:     "A",
			history: []string{"A"},
			store: func() *mockStore {
				m := &mockStore{}

				m.On("Get", "A").Return("automaton", nil).Once()

				return m
			}(),
			size:    4,
			want:    "automaton",
			wantErr: assert.NoError,
		},
		{
			name:    "hit failure",
			key:     "A",
			history: []string{"A"},
			store: func() *mockStore {
				m := &mockStore{}

				m.On("Get", "A").Return("", errors.New("failure")).Twice()

				return m
			}(),
			size: 4,
			want: "",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "backing store error: failure")
			},
		},
		{
			name:    "replacement",
			key:     "A",
			history: []string{"A", "B", "C", "D"},
			store: func() *mockStore {
				m := &mockStore{}

				m.On("Get", "A").Return("automaton", nil).Twice()
				m.On("Get", "B").Return("binary", nil).Once()
				m.On("Get", "C").Return("cache", nil).Once()
				m.On("Get", "D").Return("data", nil).Once()

				return m
			}(),
			size:    3,
			want:    "automaton",
			wantErr: assert.NoError,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c := cache.NewLRUCache[string, string](tc.store, tc.size)

			for _, key := range tc.history {
				_, _ = c.Get(key)
			}

			got, err := c.Get(tc.key)

			if tc.wantErr(t, err) {
				assert.Equal(t, tc.want, got)

				tc.store.AssertExpectations(t)
			}
		})
	}
}
