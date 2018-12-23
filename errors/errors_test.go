package errors

import (
	"errors"
	"fmt"
)

func ExampleAggregateError() {
	err := NewAggregateError("the following errors occurred:", errors.New("foo"), errors.New("bar"))
	fmt.Println(err)
	// Output:
	// the following errors occurred:
	// foo
	// bar
}

func ExampleAggregateError_withoutHeader() {
	err := NewAggregateError("", errors.New("foo"), errors.New("bar"))
	fmt.Println(err)
	// Output:
	// foo
	// bar
}

func ExampleAggregateError_headerOnly() {
	err := NewAggregateError("why would you use an aggregate error in this case?")
	fmt.Println(err)
	// Output:
	// why would you use an aggregate error in this case?
}

func ExampleErrorString_Error() {
	var err ErrorString = "I'm sorry, Dave. I'm afraid I can't do that."
	fmt.Println(err.Error())
	// Output:
	// I'm sorry, Dave. I'm afraid I can't do that.
}
