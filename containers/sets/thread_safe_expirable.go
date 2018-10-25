package sets

import (
	"sync"
	"time"
)

type threadSafeExpirableSet struct {
	mutex sync.RWMutex
	Set
}

// NewThreadSafeExpirableSet creates a new thread-safe set using a read-write
// mutex.
// The initial capacity does not bound the set's size: sets grow to accommodate
// the number of elements to store.
// TTL (time to live) specifies the duration after which an element expires.
func NewThreadSafeExpirableSet(initialCapacity int, ttl time.Duration) Set {
	return &threadSafeExpirableSet{Set: NewExpirableSet(initialCapacity, ttl)}
}

func (s *threadSafeExpirableSet) Add(element interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Add(element)
}

func (s *threadSafeExpirableSet) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Clear()
}

func (s *threadSafeExpirableSet) Contains(element interface{}) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.Contains(element)
}

func (s *threadSafeExpirableSet) Remove(element interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Remove(element)
}

func (s *threadSafeExpirableSet) Size() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.Size()
}

func (s *threadSafeExpirableSet) ToSlice() []interface{} {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.ToSlice()
}

func (s *threadSafeExpirableSet) ToStringSlice() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.ToStringSlice()
}
