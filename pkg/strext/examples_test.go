package strext

import (
	"fmt"
	"reflect"
)

func ExampleTrim() {
	var values = []string{
		"",
		" one  ",
		"two  ",
		" Three",
		"",
		"Four",
	}
	var expectedValues = []string{"one", "two", "Three", "Four"}

	newValues := Trim(values)

	if reflect.DeepEqual(newValues, expectedValues) {
		fmt.Println("String array trimmed!")
	}
	// Output: String array trimmed!
}
