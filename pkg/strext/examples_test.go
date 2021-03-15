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
	var expectedValues = []string{"", "one", "two", "Three", "", "Four"}
	var expectedTrimmedValues = []string{"one", "two", "Three", "Four"}

	newValues := TrimArray(values, true)

	if reflect.DeepEqual(newValues, expectedTrimmedValues) {
		fmt.Println("String array trimmed and empties removed!")
	}

	newValues = TrimArray(values, false)
	if reflect.DeepEqual(newValues, expectedValues) {
		fmt.Println("String array trimmed!")
	}

	// Output: String array trimmed and empties removed!
	// String array trimmed!
}

func ExampleTrimLines() {
	var value = `
This is line 1.		
	This is line 2.


This is line 5.
	`
	var expectedValue = `
This is line 1.
This is line 2.


This is line 5.
`
	trimmedLines := TrimLines(value)

	if trimmedLines == expectedValue {
		fmt.Printf("Lines trimmed!")
	}
	// Output: Lines trimmed!
}
