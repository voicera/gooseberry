package maps_test

import (
	"testing"

	"github.com/voicera/gooseberry/containers/maps"
	"github.com/voicera/tester/assert"
)

func TestGetOrDefault_keyFound(t *testing.T) {
	m := maps.NewMap(map[interface{}]interface{}{"key": "value"})
	value := m.GetOrDefault("key", "default")
	assert.For(t).ThatActual(value).Equals("value")
}

func TestGetOrDefault_keyNotFound(t *testing.T) {
	m := maps.NewMap(map[interface{}]interface{}{"key": "value"})
	value := m.GetOrDefault("not found", "default")
	assert.For(t).ThatActual(value).Equals("default")
}

func TestImmutableMapGet_changeOriginal(t *testing.T) {
	original := map[interface{}]interface{}{"name": "original"}
	m := maps.NewImmutableMap(original)
	original["name"] = "changed"
	assert.For(t).ThatActual(m.GetOrDefault("name", nil)).Equals("original")
}
