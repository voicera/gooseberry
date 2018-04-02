package queues_test

import (
	"math/rand"
	"testing"

	"github.com/voicera/gooseberry/containers/queues"
	"github.com/voicera/tester/assert"
)

func TestPriorityQueue(t *testing.T) {
	pq := queues.NewPriorityQueue(4)
	expectedValuesInOrder := []interface{}{nil, "high", "med", "low"}
	expectedPriorities := []int{10, 1, 0, -1}
	for i := range rand.Perm(len(expectedValuesInOrder)) {
		pq.Add(expectedValuesInOrder[i], expectedPriorities[i])
	}

	for i, expected := range expectedValuesInOrder {
		assert.For(t, expected).ThatActual(pq.Len()).Equals(len(expectedValuesInOrder) - i)
		value, found := pq.Poll()
		assert.For(t, expected).ThatActual(value).Equals(expected)
		assert.For(t, expected).ThatActual(found).IsTrue()
		assert.For(t, expected).ThatActual(pq.Len()).Equals(len(expectedValuesInOrder) - i - 1)
	}

	value, found := pq.Poll()
	assert.For(t, "empty").ThatActual(value).Equals(nil)
	assert.For(t, "empty").ThatActual(found).IsFalse()
	assert.For(t, "empty").ThatActual(pq.Len()).Equals(0)
}
