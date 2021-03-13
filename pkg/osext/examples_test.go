package osext

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

func ExampleFileExists() {
	if FileExists("/etc/hosts") {
		fmt.Println("File exists")
	}
	// Output: File exists
}
func ExampleFileOrDirectoryExists() {
	if FileOrDirectoryExists("/etc/hosts") {
		fmt.Println("File or directory exists")
	}
	// Output: File or directory exists
}
func ExampleIsDirectory() {
	var file = "/etc"
	if IsDirectory(file) {
		fmt.Printf("File '%s' is directory", file)
	}
	// Output: File '/etc' is directory
}

func ExampleCreateTempWithContents() {
	// Create a temp file with contents
	tempFile, err := CreateTempWithContents("", "testFile*.txt", []byte("File contents"), 0644)
	if err != nil {
		fmt.Printf("ERROR while creating temp file: %v\n", err.Error())
		return
	}
	defer os.Remove(tempFile.Name())

	fmt.Println("File created!")
	// Output: File created!
}

func Example() {
	contents := `

This is a file.

`
	contentsTrimmed := strings.TrimSpace(contents)

	// Create a temp file with contents
	tempFile, err := CreateTempWithContents("", "testFile*.txt", []byte(contents), 0644)
	if err != nil {
		fmt.Printf("ERROR while creating temp file: %v\n", err.Error())
		return
	}
	defer os.Remove(tempFile.Name())

	// Check if file exists
	if FileExists(tempFile.Name()) {
		fmt.Println("File created!")
	}

	// Read the file contents:
	var bytes []byte
	if bytes, err = ReadFile(tempFile.Name(), false); err != nil {
		fmt.Printf("ERROR while reading the file as bytes: %v\n", err.Error())
		return
	}
	if reflect.DeepEqual(bytes, []byte(contents)) {
		fmt.Println("File content as bytes read as expected!")
	}

	// Read the file contents as a trimmed string
	var readStr string
	if readStr, err = ReadFileAsString(tempFile.Name(), true); err != nil {
		fmt.Printf("ERROR while reading the file as string: %v\n", err.Error())
		return
	}
	if readStr == contentsTrimmed {
		fmt.Println("File contents as trimmed string read as expected!")
	}
	// Output: File created!
	// File content as bytes read as expected!
	// File contents as trimmed string read as expected!
}
