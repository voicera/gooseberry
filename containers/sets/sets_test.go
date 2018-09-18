package sets_test

import (
	"testing"

	"github.com/voicera/gooseberry/containers/sets"
	"github.com/voicera/tester/assert"
)

func TestSet(t *testing.T) {
	testScenarios := []func(sets.Set){add, addTwice, remove}
	for _, testScenario := range testScenarios {
		set := sets.NewSet(1)
		assert.For(t).ThatActual(set.Size()).Equals(0)
		assert.For(t).ThatActual(set.Contains(0)).IsFalse()
		assert.For(t).ThatActual(set.Contains(1)).IsFalse()
		testScenario(set)
		assert.For(t).ThatActual(set.Size()).Equals(1)
		assert.For(t).ThatActual(set.Contains(0)).IsFalse()
		assert.For(t).ThatActual(set.Contains(1)).IsTrue()
		assert.For(t).ThatActual(set.ToSlice()).Equals([]interface{}{1})
	}
}

func TestNewSetFromStrings(t *testing.T) {
	set := sets.NewSetFromStrings("foo", "bar")
	assert.For(t).ThatActual(set.Size()).Equals(2)
	assert.For(t).ThatActual(set.Contains("foo")).IsTrue()
	assert.For(t).ThatActual(set.Contains("bar")).IsTrue()
	assert.For(t).ThatActual(set.Contains("Bar")).IsFalse()
	assert.For(t).ThatActual(set.Contains("baz")).IsFalse()
}

func TestClear(t *testing.T) {
	set := sets.NewSetFromStrings("foo", "bar")
	assert.For(t).ThatActual(set.Size()).Equals(2)
	assert.For(t).ThatActual(set.Contains("foo")).IsTrue()
	assert.For(t).ThatActual(set.Contains("bar")).IsTrue()

	// Verify clearing the set removes all elements, leaving a functional empty set
	set.Clear()
	assert.For(t).ThatActual(set.Size()).Equals(0)
	assert.For(t).ThatActual(set.Contains("foo")).IsFalse()
	assert.For(t).ThatActual(set.Contains("bar")).IsFalse()

	// Verify the set is still functional and can be added to
	set.Add("foo")
	assert.For(t).ThatActual(set.Size()).Equals(1)
	assert.For(t).ThatActual(set.Contains("foo")).IsTrue()

	// Verify clearing the set after additional adds
	set.Clear()
	assert.For(t).ThatActual(set.Size()).Equals(0)
	assert.For(t).ThatActual(set.Contains("foo")).IsFalse()
}

func add(set sets.Set) {
	set.Add(1)
}

func addTwice(set sets.Set) {
	set.Add(1)
	set.Add(1)
}

func remove(set sets.Set) {
	set.Remove(1)
	set.Add(0)
	set.Add(1)
	set.Remove(0)
}
