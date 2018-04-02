package sets

import (
	"time"
)

type expirableSet struct {
	underlying map[interface{}]time.Time
	ttl        time.Duration
}

// NewExpirableSet creates a new set. The initial capacity does not bound
// the set's size: sets grow to accommodate the number of elements to store.
// TTL (time to live) specifies the duration after which an element expires.
// TODO(Geish): shrink the size on expiry
func NewExpirableSet(initialCapacity int, ttl time.Duration) Set {
	return &expirableSet{underlying: make(map[interface{}]time.Time, initialCapacity), ttl: ttl}
}

func (s *expirableSet) Add(element interface{}) {
	s.underlying[element] = time.Now().UTC()
}

func (s *expirableSet) Contains(element interface{}) bool {
	timestamp, found := s.underlying[element]
	return found && time.Now().UTC().Before(timestamp.Add(s.ttl))
}

func (s *expirableSet) Remove(element interface{}) {
	delete(s.underlying, element)
}

func (s *expirableSet) Size() int {
	return len(s.underlying)
}

func (s *expirableSet) ToSlice() []interface{} {
	slice := make([]interface{}, 0, s.Size())
	for element, timestamp := range s.underlying {
		if time.Now().UTC().Before(timestamp.Add(s.ttl)) {
			slice = append(slice, element)
		}
	}
	return slice
}

func (s *expirableSet) ToStringSlice() []string {
	slice := make([]string, 0, s.Size())
	for element, timestamp := range s.underlying {
		if time.Now().UTC().Before(timestamp.Add(s.ttl)) {
			slice = append(slice, element.(string)) // TODO(Geish): check type and call fmt.Sprint if not a string?
		}
	}
	return slice
}
