// Package maps provides map data structures.
package maps

// Map is a map data structure.
type Map interface {
	// Get returns the value corresponding to the specified key if found.
	// It also returns a flag to indicate whether the key exists.
	Get(key interface{}) (interface{}, bool)

	// GetOrDefault returns the value to corresponding to the specified key,
	// or the specified defaultValue if no such value exists.
	GetOrDefault(key interface{}, defaultValue interface{}) interface{}
}

type mapDecorator struct {
	underlying map[interface{}]interface{}
}

// NewImmutableMap creates a new immutable map.
func NewImmutableMap(underlying map[interface{}]interface{}) Map {
	immutableCopy := make(map[interface{}]interface{}, len(underlying))
	for k, v := range underlying {
		immutableCopy[k] = v
	}
	return NewMap(immutableCopy)
}

// NewMap creates a new map.
func NewMap(underlying map[interface{}]interface{}) Map {
	return &mapDecorator{underlying}
}

func (m *mapDecorator) Get(key interface{}) (interface{}, bool) {
	value, found := m.underlying[key]
	return value, found
}

func (m *mapDecorator) GetOrDefault(key interface{}, defaultValue interface{}) interface{} {
	value, found := m.underlying[key]
	if !found {
		value = defaultValue
	}
	return value
}
