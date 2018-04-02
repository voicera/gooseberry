/*
Package sets provides set data structures.
The set data structures in this package do not replace the recommended
implementation in Effective Go:
A set can be implemented as a map with value type bool. Set the map entry to
true to put the value in the set, and then test it by simple indexing.
	attended := map[string]bool{
		"Ann": true,
		"Joe": true,
		...
	}
	if attended[person] { // will be false if person is not in the map
		fmt.Println(person, "was at the meeting")
	}
However, this package provide set types that are more efficient for
Add, Remove, and Contains calls; and more functionality on top of that.
Due to the lack of generics in Go, it's recommended to choose
the map[string]bool when the set must be iterated many times.
*/
package sets
