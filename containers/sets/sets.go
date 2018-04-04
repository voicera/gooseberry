package sets

// Set is a set data structure.
// It does not replace the recommended implementation in Effective Go:
// A set can be implemented as a map with value type bool. Set the map entry to
// true to put the value in the set, and then test it by simple indexing.
//
// 	attended := map[string]bool{
// 		"Ann": true,
// 		"Joe": true,
// 		...
// 	}
//
// 	if attended[person] { // will be false if person is not in the map
// 		fmt.Println(person, "was at the meeting")
// 	}
//
// However, this set type is more efficient for Add, Remove, and Contains calls.
// Due to the lack of generics in Go, it's recommended to choose
// the map[string]bool when the set must be iterated many times.
type Set interface {
	// Add adds the specified element to the set.
	Add(element interface{})

	// Contains checks whether or not the specified element belongs to the set.
	Contains(element interface{}) bool

	// Remove removes the specified element from the set.
	Remove(element interface{})

	// Size return the number of elements in the set.
	Size() int

	// ToSlice converts the set into a slice.
	ToSlice() []interface{}

	// ToStringSlice converts the set into a string slice.
	ToStringSlice() []string
}

type set map[interface{}]struct{}

var present = struct{}{} // 0 bytes

// NewSet creates a new set. The initial capacity does not bound the set's size:
// sets grow to accommodate the number of elements to store.
func NewSet(initialCapacity int) Set {
	return make(set, initialCapacity)
}

// NewSetFromStrings creates a new set from the specified strings.
func NewSetFromStrings(strings ...string) Set {
	set := NewSet(len(strings))
	for _, s := range strings {
		set.Add(s)
	}
	return set
}

func (s set) Add(element interface{}) {
	s[element] = present
}

func (s set) Contains(element interface{}) bool {
	_, found := s[element]
	return found
}

func (s set) Remove(element interface{}) {
	delete(s, element)
}

func (s set) Size() int {
	return len(s)
}

func (s set) ToSlice() []interface{} {
	slice := make([]interface{}, 0, s.Size())
	for element := range s {
		slice = append(slice, element)
	}
	return slice
}

func (s set) ToStringSlice() []string {
	slice := make([]string, 0, s.Size())
	for element := range s {
		slice = append(slice, element.(string)) // TODO (Geish): check type and call fmt.Sprint if not a string?
	}
	return slice
}
